'use client';
import { useState } from 'react';
import { useParams } from 'next/navigation';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { Heart } from 'lucide-react';
import api from '@/lib/api';
import MovieCard from '@/components/movie-card';
import BookCard from '@/components/book-card';
import { Bookmark } from 'lucide-react';
import { useBookmark } from '@/hooks/useBookmark';
import { useLike } from '@/hooks/useLike';

export default function BookPage() {
  const params = useParams();
  const bookId = params.id as string;
  const queryClient = useQueryClient();
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

  const { data: similarBooks } = useQuery<any[]>({
    queryKey: ['similarBooks', bookId],
    queryFn: async () => {
      const res = await api.post('/api/recommend', {
        selected_ids: [`book_${bookId}`],
        direction: 'book_to_book',
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
    enabled: !!bookId,
  });

  const { data: similarMovies } = useQuery<any[]>({
    queryKey: ['similarMovies', bookId],
    queryFn: async () => {
      const res = await api.post('/api/recommend', {
        selected_ids: [`book_${bookId}`],
        direction: 'book_to_movie',
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
    enabled: !!bookId,
  });

  if (bookLoading) return <div className="p-8 text-center">Loading...</div>;
  if (!book) return <div className="p-8 text-center">Book not found</div>;

  return (
    <div className="min-h-screen pb-20 px-4 max-w-screen-md mx-auto overflow-x-hidden">
      <div className="pt-6">
        {book.image && (
          <img
            src={book.image}
            alt={book.title}
            className="w-full max-w-full h-48 sm:h-64 object-cover rounded-xl mb-4"
          />
        )}
        <div className="flex items-start justify-between gap-2">
          <h1 className="text-2xl font-bold">{book.title}</h1>
          <div className="flex gap-2">
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
        {book.authors && <p className="text-muted-foreground">by {book.authors}</p>}
        {book.genres && (
            <div className="flex flex-wrap gap-1 mt-2">
                {normalizeGenres(book.genres).map((g: string) => (
                    <span key={g} className="px-2 py-0.5 bg-secondary rounded-full text-xs">{g}</span>
                ))}
            </div>
        )}
        <div className="flex items-center gap-4 mt-2">
          {book.rating > 0 && <p>★ {book.rating}</p>}
          {book.year > 0 && <p className="text-muted-foreground">{book.year}</p>}
        </div>
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
      </div>

      {similarBooks && similarBooks.length > 0 && (
        <div className="mt-8">
          <h2 className="text-xl font-semibold mb-4">Similar Books</h2>
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

      {similarMovies && similarMovies.length > 0 && (
        <div className="mt-8">
          <h2 className="text-xl font-semibold mb-4">Movies You Might Like</h2>
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
    </div>
  );
}