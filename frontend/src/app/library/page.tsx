'use client';

import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import api from '@/lib/api';
import { Book, Movie } from '@/types';
import { Button } from '@/components/ui/button';
import { loadBooks, loadMovies } from '@/lib/data';
import { useEffect, useState } from 'react';

export default function LibraryPage() {
  const queryClient = useQueryClient();

  const { data: likedBooks } = useQuery<string[]>({
    queryKey: ['likedBooks'],
    queryFn: async () => {
      const res = await api.get('/api/user/library/books');
      return res.data.books.map((b: any) => b.book_id);
    },
  });

  const { data: likedMovies } = useQuery<string[]>({
    queryKey: ['likedMovies'],
    queryFn: async () => {
      const res = await api.get('/api/user/library/movies');
      return res.data.movies.map((m: any) => m.movie_id);
    },
  });

  const removeBook = useMutation({
    mutationFn: (bookId: string) => api.delete(`/api/book/${bookId}/like`),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['likedBooks'] }),
  });

  const removeMovie = useMutation({
    mutationFn: (movieId: string) => api.delete(`/api/movie/${movieId}/like`),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['likedMovies'] }),
  });

  return (
    <div className="container mx-auto p-4">
      <h1 className="text-2xl font-bold mb-4">Your Library</h1>
      <h2 className="text-xl font-semibold mb-2">Books</h2>
      <div className="space-y-2">
        {likedBooks?.map((id) => (
          <div key={id} className="flex justify-between items-center border p-2 rounded">
            <span>{id} <BookTitle bookId={id} /></span>
            <Button size="sm" variant="destructive" onClick={() => removeBook.mutate(id)}>Remove</Button>
          </div>
        ))}
      </div>
      <h2 className="text-xl font-semibold mt-6 mb-2">Movies</h2>
      <div className="space-y-2">
        {likedMovies?.map((id) => (
          <div key={id} className="flex justify-between items-center border p-2 rounded">
            <span>{id} <MovieTitle movieId={id} /></span>
            <Button size="sm" variant="destructive" onClick={() => removeMovie.mutate(id)}>Remove</Button>
          </div>
        ))}
      </div>
    </div>
  );
}

// Вспомогательные компоненты для отображения названий
function BookTitle({ bookId }: { bookId: string }) {
  const [title, setTitle] = useState('');
  useEffect(() => {
    loadBooks().then((books) => {
      const b = books.find((b) => b.book_id === bookId);
      if (b) setTitle(b.title);
    });
  }, [bookId]);
  return <span className="ml-2 text-sm text-muted-foreground">{title}</span>;
}

function MovieTitle({ movieId }: { movieId: string }) {
  const [title, setTitle] = useState('');
  useEffect(() => {
    loadMovies().then((movies) => {
      const m = movies.find((m) => m.movie_id === movieId);
      if (m) setTitle(m.title);
    });
  }, [movieId]);
  return <span className="ml-2 text-sm text-muted-foreground">{title}</span>;
}