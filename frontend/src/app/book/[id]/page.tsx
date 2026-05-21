'use client';
import { useEffect, useState } from 'react';
import { useParams } from 'next/navigation';
import { useQuery } from '@tanstack/react-query';
import api from '@/lib/api';
import MovieCard from '@/components/movie-card';
import BookCard from '@/components/book-card';

export default function BookPage() {
  const params = useParams();
  const bookId = params.id as string;
  const [book, setBook] = useState<any>(null);

  // Загружаем информацию о книге
  useEffect(() => {
    api.get(`/api/items/batch?ids=${bookId}&type=book`)
      .then(res => {
        const items = res.data.items || [];
        if (items.length > 0) setBook(items[0]);
      })
      .catch(() => setBook(null));
  }, [bookId]);

  // Похожие книги
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
      // Получаем метаданные через batch
      const batch = await api.get(`/api/items/batch?ids=${ids.join(',')}&type=book`);
      return batch.data.items || [];
    },
    enabled: !!bookId,
  });

  // Похожие фильмы
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

  if (!book) return <div className="p-8 text-center">Loading...</div>;

  return (
    <div className="min-h-screen pb-20 px-4">
      <div className="pt-6">
        {book.image_url && (
          <img src={book.image_url} alt={book.title} className="w-full h-64 object-cover rounded-xl mb-4" />
        )}
        <h1 className="text-2xl font-bold">{book.title}</h1>
        {book.authors && <p className="text-muted-foreground">by {book.authors}</p>}
        {book.genres && (
          <div className="flex flex-wrap gap-1 mt-2">
            {book.genres.map((g: string) => (
              <span key={g} className="px-2 py-0.5 bg-secondary rounded-full text-xs">{g}</span>
            ))}
          </div>
        )}
        <div className="flex items-center gap-4 mt-2">
          {book.rating > 0 && <p>★ {book.rating}</p>}
          {book.year > 0 && <p className="text-muted-foreground">{book.year}</p>}
        </div>
        {book.description && <p className="mt-4 text-sm">{book.description}</p>}
      </div>

      {similarBooks && similarBooks.length > 0 && (
        <div className="mt-8">
          <h2 className="text-xl font-semibold mb-4">Similar Books</h2>
          <div className="flex gap-2 overflow-x-auto pb-2 no-scrollbar">
            {similarBooks.slice(0, 10).map((item: any) => (
              <div key={item.id} className="flex-shrink-0 w-28">
                <BookCard book={{ book_id: item.id, title: item.title, image_url: item.image }} />
              </div>
            ))}
          </div>
        </div>
      )}

      {similarMovies && similarMovies.length > 0 && (
        <div className="mt-8">
          <h2 className="text-xl font-semibold mb-4">Movies You Might Like</h2>
          <div className="flex gap-2 overflow-x-auto pb-2 no-scrollbar">
            {similarMovies.slice(0, 10).map((item: any) => (
              <div key={item.id} className="flex-shrink-0 w-28">
                <MovieCard movie={{ movie_id: item.id, title: item.title, poster_url: item.image }} />
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  );
}