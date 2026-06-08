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

export default function BookPage() {
  const params = useParams();
  const bookId = params.id as string;
  const router = useRouter();
  const [descExpanded, setDescExpanded] = useState(false);

  const { isLiked, toggleLike, isPending: likePending } = useLike('book', bookId);
  const { isBookmarked, toggleBookmark, isPending: bookmarkPending } = useBookmark('book', bookId);

  const { data: book, isLoading: bookLoading } = useQuery<any>({
    queryKey: ['item', bookId, 'book'],
    queryFn: async () => {
      const res = await api.get(`/api/items/batch?ids=${bookId}&type=book`);
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

  const { data: similarBooks, isLoading: booksLoading } = useQuery<any[]>({
    queryKey: ['similarBooks', bookId],
    queryFn: async () => {
      const res = await api.post('/api/recommend', {
        selected_ids: [`book_${bookId}`],
        direction: 'book_to_book',
        weights: { genre: 0.3, text: 0.4, image: 0.3 },
        contextual: true,
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
    enabled: !!bookId,
  });

  const { data: similarMovies, isLoading: moviesLoading } = useQuery<any[]>({
    queryKey: ['similarMovies', bookId],
    queryFn: async () => {
      const res = await api.post('/api/recommend', {
        selected_ids: [`book_${bookId}`],
        direction: 'book_to_movie',
        weights: { genre: 0.3, text: 0.4, image: 0.3 },
        contextual: true,
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
    enabled: !!bookId,
  });

  if (bookLoading) return <div className="p-8 text-center">Loading...</div>;
  if (!book) return <div className="p-8 text-center">Book not found</div>;

  return (
    <div className="relative min-h-screen pb-20 overflow-x-hidden">
      {book.image && (
        <>
          <img
            src={book.image}
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
          {book.image && (
            <div className="w-1/3 flex-shrink-0">
              <img
                src={book.image}
                alt={book.title}
                className="w-full h-auto object-cover rounded-xl"
                style={{ aspectRatio: '3/4' }}
              />
            </div>
          )}
          <div className="flex-1 min-w-0 bg-background/70 backdrop-blur-sm rounded-xl p-3 shadow-md">
            <div className="flex items-start justify-between gap-2">
              <h1 className="text-2xl font-bold">{book.title}</h1>
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
            {book.authors && <p className="text-muted-foreground text-sm">by {book.authors}</p>}
            <div className="flex items-center gap-4 mt-2">
              {book.rating > 0 && <p>★ {Math.floor(Number(book.rating) * 10) / 10}</p>}
              {book.year > 0 && <p className="text-muted-foreground text-sm">{book.year}</p>}
            </div>
          </div>
        </div>

        {book.genres && (
          <div className="mt-3 overflow-x-auto no-scrollbar">
            <div className="flex gap-1 min-w-max">
              {normalizeGenres(book.genres).map((g: string) => (
                <span key={g} className="px-2 py-0.5 bg-secondary rounded-full text-xs whitespace-nowrap">{g}</span>
              ))}
            </div>
          </div>
        )}

        {book.description && (
          <div className="mt-4">
            <p className="text-sm whitespace-pre-wrap break-words">
              {descExpanded || book.description.length <= 300
                ? book.description
                : `${book.description.slice(0, 300)}...`}
            </p>
            {book.description.length > 300 && (
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
          <h2 className="text-xl font-semibold mb-4">Similar Books</h2>
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

        <div className="mt-8">
          <h2 className="text-xl font-semibold mb-4">Movies You Might Like</h2>
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
      </div>
    </div>
  );
}