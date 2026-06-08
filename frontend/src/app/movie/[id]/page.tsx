'use client';
import { useState } from 'react';
import { useParams, useRouter } from 'next/navigation';
import { useQuery } from '@tanstack/react-query';
import { Heart, Bookmark, Loader2, ArrowLeft } from 'lucide-react';
import api from '@/lib/api';
import MovieCard from '@/components/movie-card';
import BookCard from '@/components/book-card';
import { useBookmark } from '@/hooks/useBookmark';
import { useLike } from '@/hooks/useLike';

export default function MoviePage() {
  const params = useParams();
  const movieId = params.id as string;
  const router = useRouter();
  const [descExpanded, setDescExpanded] = useState(false);

  const { isLiked, toggleLike, isPending: likePending } = useLike('movie', movieId);
  const { isBookmarked, toggleBookmark, isPending: bookmarkPending } = useBookmark('movie', movieId);

  const { data: movie, isLoading: movieLoading } = useQuery<any>({
    queryKey: ['item', movieId, 'movie'],
    queryFn: async () => {
      const res = await api.get(`/api/items/batch?ids=${movieId}&type=movie`);
      return res.data.items?.[0] ?? null;
    },
    staleTime: 10 * 60 * 1000,
  });

  const normalizeGenres = (genres: any): string[] => {
    let items: string[] = [];
    if (Array.isArray(genres)) {
      items = genres.map(g => String(g));
    } else if (typeof genres === 'string') {
      items = genres.split(',').map(s => s.trim());
    }
    const flat = items.flatMap(item =>
      item.split(',').map(s => {
        const trimmed = s.trim().toLowerCase();
        return trimmed.charAt(0).toUpperCase() + trimmed.slice(1);
      })
    );
    return [...new Set(flat)].filter(Boolean);
  };

  const { data: similarMovies, isLoading: moviesLoading } = useQuery<any[]>({
    queryKey: ['similarMovies', movieId],
    queryFn: async () => {
      const res = await api.post('/api/recommend', {
        selected_ids: [`movie_${movieId}`],
        direction: 'movie_to_movie',
        weights: { genre: 0.3, text: 0.4, image: 0.3 },
      });
      let result = await api.get(`/api/result/${res.data.task_id}`);
      while (result.data.status === 'pending') {
        await new Promise((r) => setTimeout(r, 1000));
        result = await api.get(`/api/result/${res.data.task_id}`);
      }
      const ids = result.data.movies || [];
      if (ids.length === 0) return [];
      const batch = await api.get(`/api/items/batch?ids=${ids.join(',')}&type=movie`);
      return batch.data.items || [];
    },
    enabled: !!movieId,
  });

  const { data: similarBooks, isLoading: booksLoading } = useQuery<any[]>({
    queryKey: ['similarBooks', movieId],
    queryFn: async () => {
      const res = await api.post('/api/recommend', {
        selected_ids: [`movie_${movieId}`],
        direction: 'movie_to_book',
        weights: { genre: 0.3, text: 0.4, image: 0.3 },
      });
      let result = await api.get(`/api/result/${res.data.task_id}`);
      while (result.data.status === 'pending') {
        await new Promise((r) => setTimeout(r, 1000));
        result = await api.get(`/api/result/${res.data.task_id}`);
      }
      const ids = result.data.movies || [];
      if (ids.length === 0) return [];
      const batch = await api.get(`/api/items/batch?ids=${ids.join(',')}&type=book`);
      return batch.data.items || [];
    },
    enabled: !!movieId,
  });

  if (movieLoading) return <div className="p-8 text-center">Loading...</div>;
  if (!movie) return <div className="p-8 text-center">Movie not found</div>;

  const castList = movie.cast?.split(',').map((s: string) => s.trim()).filter(Boolean) || [];
  const displayedCast = castList.slice(0, 3).join(', ');
  const extraCount = castList.length - 3;

  return (
    <div className="relative min-h-screen pb-20 overflow-x-hidden">
      {movie.image && (
        <>
          <img
            src={movie.image}
            alt=""
            className="absolute inset-0 w-full h-[35vh] object-cover z-0"
            style={{
              maskImage: 'linear-gradient(to bottom, black 20%, transparent 100%)',
              WebkitMaskImage: 'linear-gradient(to bottom, black 20%, transparent 100%)',
            }}
          />
          <div className="absolute inset-0 bg-background/35 z-[1] h-[35vh]" />
        </>
      )}

      <div className="relative z-10 px-4 max-w-screen-md mx-auto">
        <div className="pt-6 pb-2">
          <button
            onClick={() => router.back()}
            className="p-2 rounded-full bg-background/60 backdrop-blur-sm shadow-md hover:bg-background/80 transition"
            aria-label="Go back"
          >
            <ArrowLeft className="w-5 h-5" />
          </button>
        </div>

        <div className="flex gap-4 items-center">
          {movie.image && (
            <div className="w-1/3 flex-shrink-0">
              <img
                src={movie.image}
                alt={movie.title}
                className="w-full h-auto object-cover rounded-xl"
                style={{ aspectRatio: '2/3' }}
              />
            </div>
          )}
          <div className="flex-1 min-w-0 bg-background/70 backdrop-blur-sm rounded-xl p-3 shadow-md">
            <div className="flex items-start justify-between gap-2">
              <h1 className="text-2xl font-bold">{movie.title}</h1>
              <div className="flex gap-2 flex-shrink-0">
                <button onClick={toggleBookmark} disabled={bookmarkPending} className={`p-2 rounded-full transition ${isBookmarked ? 'bg-blue-100 text-blue-500' : 'bg-secondary/50 text-muted-foreground hover:bg-blue-100 hover:text-blue-500'} ${bookmarkPending ? 'opacity-50 pointer-events-none' : ''}`}>
                  <Bookmark className={`w-5 h-5 ${isBookmarked ? 'fill-current' : ''}`} />
                </button>
                <button
                  onClick={(e) => {
                    e.preventDefault();
                    toggleLike();
                  }}
                  disabled={likePending}
                  className={`p-2 rounded-full transition ${isLiked ? 'bg-red-100 text-red-500' : 'bg-secondary/50 text-muted-foreground hover:bg-red-100 hover:text-red-500'} ${likePending ? 'opacity-50 pointer-events-none' : ''}`}
                  title={isLiked ? 'Remove from library' : 'Add to library'}
                >
                  <Heart className={`w-5 h-5 ${isLiked ? 'fill-current' : ''}`} />
                </button>
              </div>
            </div>
            {movie.director && <p className="text-muted-foreground text-sm">Director: {movie.director}</p>}
            {movie.cast && (
              <p className="text-sm text-muted-foreground">
                Cast: {displayedCast}{extraCount > 0 ? ` and ${extraCount} more` : ''}
              </p>
            )}
            <div className="flex items-center gap-4 mt-2">
              {movie.rating > 0 && <p>★ {Math.floor(Number(movie.rating) * 10) / 10}</p>}
              {movie.year > 0 && <p className="text-muted-foreground text-sm">{movie.year}</p>}
              {movie.runtime > 0 && <p className="text-muted-foreground text-sm">{movie.runtime} min</p>}
            </div>
          </div>
        </div>

        {movie.genres && (
          <div className="mt-3 overflow-x-auto no-scrollbar">
            <div className="flex gap-1 min-w-max">
              {normalizeGenres(movie.genres).map((g: string) => (
                <span key={g} className="px-2 py-0.5 bg-secondary rounded-full text-xs whitespace-nowrap">{g}</span>
              ))}
            </div>
          </div>
        )}

        {movie.description && (
          <div className="mt-4">
            <p className="text-sm whitespace-pre-wrap break-words">
              {descExpanded || movie.description.length <= 300
                ? movie.description
                : `${movie.description.slice(0, 300)}...`}
            </p>
            {movie.description.length > 300 && (
              <button
                onClick={() => setDescExpanded(!descExpanded)}
                className="text-primary text-sm mt-1 hover:underline"
              >
                {descExpanded ? 'Show less' : 'Show more'}
              </button>
            )}
          </div>
        )}

        <div className="mt-8">
          <h2 className="text-xl font-semibold mb-4">Similar Movies</h2>
          {moviesLoading ? (
            <div className="flex justify-center py-8">
              <Loader2 className="w-6 h-6 animate-spin text-muted-foreground" />
            </div>
          ) : similarMovies && similarMovies.length > 0 ? (
            <div className="scroll-container pb-2">
              <div className="flex gap-2 min-w-max">
                {similarMovies.slice(0, 10).map((item: any) => (
                  <div key={item.id} className="w-28 flex-shrink-0">
                    <MovieCard movie={{ movie_id: item.id, title: item.title, poster_url: item.image }} />
                  </div>
                ))}
              </div>
            </div>
          ) : null}
        </div>

        <div className="mt-8">
          <h2 className="text-xl font-semibold mb-4">Books You Might Like</h2>
          {booksLoading ? (
            <div className="flex justify-center py-8">
              <Loader2 className="w-6 h-6 animate-spin text-muted-foreground" />
            </div>
          ) : similarBooks && similarBooks.length > 0 ? (
            <div className="scroll-container pb-2">
              <div className="flex gap-2 min-w-max">
                {similarBooks.slice(0, 10).map((item: any) => (
                  <div key={item.id} className="w-28 flex-shrink-0">
                    <BookCard book={{ book_id: item.id, title: item.title, image_url: item.image }} />
                  </div>
                ))}
              </div>
            </div>
          ) : null}
        </div>
      </div>
    </div>
  );
}