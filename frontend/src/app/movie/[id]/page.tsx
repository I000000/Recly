'use client';
import { useState } from 'react';
import { useParams } from 'next/navigation';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { Heart } from 'lucide-react';
import api from '@/lib/api';
import MovieCard from '@/components/movie-card';
import BookCard from '@/components/book-card';

export default function MoviePage() {
  const params = useParams();
  const movieId = params.id as string;
  const queryClient = useQueryClient();
  const [descExpanded, setDescExpanded] = useState(false);

  const { data: movie, isLoading: movieLoading } = useQuery<any>({
    queryKey: ['item', movieId, 'movie'],
    queryFn: async () => {
      const res = await api.get(`/api/items/batch?ids=${movieId}&type=movie`);
      return res.data.items?.[0] ?? null;
    },
    staleTime: 10 * 60 * 1000,
  });

  const { data: likedMovies = [] } = useQuery<string[]>({
    queryKey: ['likedMovies'],
    queryFn: async () => (await api.get('/api/user/library/movies')).data.movies.map((m: any) => m.movie_id),
    staleTime: 1000 * 60 * 30,
  });

  const isLiked = likedMovies.includes(movieId);

  const likeMovie = useMutation({
    mutationFn: () => api.post(`/api/movie/${movieId}/like`),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['likedMovies'] });
    },
  });

  const { data: similarMovies } = useQuery<any[]>({
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

  const { data: similarBooks } = useQuery<any[]>({
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
  const displayedCast = castList.slice(0, 5).join(', ');
  const extraCount = castList.length - 5;

  return (
    <div className="min-h-screen pb-20 px-4 max-w-screen-md">
      <div className="pt-6">
        {movie.image && (
          <img src={movie.image} alt={movie.title} className="w-full h-64 object-cover rounded-xl mb-4" />
        )}
        <div className="flex items-start justify-between gap-2">
          <h1 className="text-2xl font-bold">{movie.title}</h1>
          <button
            onClick={() => !isLiked && likeMovie.mutate()}
            disabled={likeMovie.isPending || isLiked}
            className={`p-2 rounded-full ${
              isLiked ? 'bg-red-100 text-red-500' : 'bg-secondary/50 text-muted-foreground hover:bg-red-100 hover:text-red-500'
            } transition`}
            title={isLiked ? 'In your library' : 'Add to library'}
          >
            <Heart className={`w-5 h-5 ${isLiked ? 'fill-current' : ''}`} />
          </button>
        </div>
        {movie.director && <p className="text-muted-foreground">Director: {movie.director}</p>}
        {movie.cast && (
          <p className="text-sm text-muted-foreground">
            Cast: {displayedCast}{extraCount > 0 ? ` and ${extraCount} more` : ''}
          </p>
        )}
        {movie.genres && (
          <div className="flex flex-wrap gap-1 mt-2">
            {movie.genres.map((g: string) => (
              <span key={g} className="px-2 py-0.5 bg-secondary rounded-full text-xs">{g}</span>
            ))}
          </div>
        )}
        <div className="flex items-center gap-4 mt-2">
          {movie.rating > 0 && <p>★ {movie.rating}</p>}
          {movie.year > 0 && <p className="text-muted-foreground">{movie.year}</p>}
          {movie.runtime > 0 && <p className="text-muted-foreground">{movie.runtime} min</p>}
        </div>
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
      </div>

      {similarMovies && similarMovies.length > 0 && (
        <div className="mt-8">
          <h2 className="text-xl font-semibold mb-4">Similar Movies</h2>
          <div className="scroll-container pb-2">
            <div className="flex gap-2 min-w-max">
              {similarMovies.slice(0, 10).map((item: any) => (
                <div key={item.id} className="w-28 flex-shrink-0">
                  <MovieCard movie={{ movie_id: item.id, title: item.title, poster_url: item.image }} />
                </div>
              ))}
            </div>
          </div>
        </div>
      )}

      {similarBooks && similarBooks.length > 0 && (
        <div className="mt-8">
          <h2 className="text-xl font-semibold mb-4">Books You Might Like</h2>
          <div className="scroll-container pb-2">
            <div className="flex gap-2 min-w-max">
              {similarBooks.slice(0, 10).map((item: any) => (
                <div key={item.id} className="w-28 flex-shrink-0">
                  <BookCard book={{ book_id: item.id, title: item.title, image_url: item.image }} />
                </div>
              ))}
            </div>
          </div>
        </div>
      )}
    </div>
  );
}