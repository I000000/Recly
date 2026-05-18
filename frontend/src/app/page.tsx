'use client';

import { useState, useEffect } from 'react';
import { useMutation, useQuery } from '@tanstack/react-query';
import api from '@/lib/api';
import { loadMovies, loadBooks } from '@/lib/data';
import { Direction, RecommendationResult, Movie, Book } from '@/types';
import DirectionSwitcher from '@/components/direction-switcher';
import MovieCard from '@/components/movie-card';
import BookCard from '@/components/book-card';

export default function HomePage() {
  const [direction, setDirection] = useState<Direction>('book_to_movie');
  const [taskId, setTaskId] = useState<string | null>(null);
  const [moviesMap, setMoviesMap] = useState<Record<string, Movie>>({});
  const [booksMap, setBooksMap] = useState<Record<string, Book>>({});

  // Загружаем метаданные один раз
  useEffect(() => {
    loadMovies().then((movies) => {
      const map: Record<string, Movie> = {};
      movies.forEach((m: any) => {
        map[String(m.movie_id)] = m;
      });
      setMoviesMap(map);
    });
    loadBooks().then((books) => {
      const map: Record<string, Book> = {};
      books.forEach((b: any) => {
        map[String(b.book_id)] = b;
      });
      setBooksMap(map);
    });
  }, []);

  const mutation = useMutation({
    mutationFn: async () => {
      const res = await api.post('/api/recommend', {
        selected_ids: [],
        direction,
        weights: { genre: 0.3, text: 0.4, image: 0.3 },
      });
      return res.data.task_id;
    },
    onSuccess: (newTaskId) => {
      setTaskId(newTaskId);
    },
  });

  useEffect(() => {
    setTaskId(null);
    mutation.mutate();
  }, [direction]);

  const { data: result } = useQuery<RecommendationResult>({
    queryKey: ['recommendationResult', taskId],
    queryFn: async () => {
      const res = await api.get(`/api/result/${taskId}`);
      return res.data;
    },
    enabled: !!taskId,
    refetchInterval: (query) => {
      if (query.state.data?.status === 'pending') return 2000;
      return false;
    },
  });

  return (
    <div className="container mx-auto p-4">
      <h1 className="text-2xl font-bold mb-4">Your Recommendations</h1>
      <DirectionSwitcher value={direction} onChange={(dir) => setDirection(dir)} />

      {mutation.isPending && <p className="mt-4">Sending request…</p>}
      {result?.status === 'pending' && <p className="mt-4">Processing your recommendations…</p>}

      {result?.status === 'done' && (
        <>
          <h2 className="text-xl font-semibold mt-4">Recommended for you</h2>
          <div className="grid grid-cols-2 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-2 mt-4">
            {result.movies?.map((id) => {
              const movie = moviesMap[String(id)];
              return movie ? <MovieCard key={id} movie={movie} /> : null;
            })}
            {result.books?.map((id) => {
              const book = booksMap[String(id)];
              return book ? <BookCard key={id} book={book} /> : null;
            })}
          </div>
        </>
      )}

      {result?.status === 'error' && <p className="mt-4 text-red-500">Something went wrong.</p>}
    </div>
  );
}