'use client';

import { useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import api from '@/lib/api';
import MovieCard from '@/components/movie-card';
import BookCard from '@/components/book-card';

type Tab = 'movies' | 'books';

async function fetchRecommendations(tab: Tab) {
  const direction = tab === 'movies' ? 'book_to_movie' : 'movie_to_book';
  // 1. Создаём задачу
  const { data: task } = await api.post('/api/recommend', {
    selected_ids: [],
    direction,
    weights: { genre: 0.3, text: 0.4, image: 0.3 },
  });
  const taskId = task.task_id;

  // 2. Опрос результата
  let result: any;
  while (true) {
    const { data: poll } = await api.get(`/api/result/${taskId}`);
    if (poll.status === 'done' || poll.status === 'error') {
      result = poll;
      break;
    }
    await new Promise(r => setTimeout(r, 1000));
  }

  if (result.status !== 'done' || !result.movies?.length) {
    throw new Error('No recommendations yet.');
  }

  // 3. Загружаем метаданные
  const ids = result.movies;
  const type = tab === 'movies' ? 'movie' : 'book';
  const { data: batch } = await api.get(`/api/items/batch?ids=${ids.join(',')}&type=${type}`);
  const items = batch.items || [];

  // Приводим к нужным полям
  return items.map((item: any) => ({
    id: item.id,
    title: item.title,
    poster_url: item.image || '',
    image_url: item.image || '',
  }));
}

export default function HomePage() {
  const [tab, setTab] = useState<Tab>('movies');

  const {
    data: cards = [],
    isLoading,
    isError,
    error,
  } = useQuery({
    queryKey: ['homeRecommendations', tab],
    queryFn: () => fetchRecommendations(tab),
    staleTime: 15 * 60 * 1000, // 15 минут данные считаются свежими
    gcTime: 30 * 60 * 1000,    // ещё 30 минут хранятся в кэше после неиспользования
    retry: 1,                  // одна повторная попытка при ошибке
    refetchOnWindowFocus: false, // не перезапрашивать при возврате в окно
  });

  return (
    <div className="min-h-screen pb-20">
      <div className="px-4 pt-6 pb-2">
        <h1 className="text-2xl font-bold tracking-tight">Recly</h1>
      </div>
      <div className="sticky top-0 z-10 bg-background/80 backdrop-blur-sm border-b border-border px-4 py-2">
        <div className="flex gap-2">
          <button
            onClick={() => setTab('movies')}
            className={`px-4 py-2 rounded-full text-sm font-medium ${tab === 'movies' ? 'bg-primary text-primary-foreground' : 'bg-secondary'}`}
          >
            Movies
          </button>
          <button
            onClick={() => setTab('books')}
            className={`px-4 py-2 rounded-full text-sm font-medium ${tab === 'books' ? 'bg-primary text-primary-foreground' : 'bg-secondary'}`}
          >
            Books
          </button>
        </div>
      </div>

      {isLoading && <div className="flex justify-center py-20">Loading...</div>}
      {isError && <p className="text-center py-20 text-destructive">{(error as any)?.message || 'Failed'}</p>}

      {cards.length > 0 && (
        <div className="grid grid-cols-2 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-2 p-4">
          {cards.map((card: any) =>
            tab === 'movies' ? (
              <MovieCard
                key={card.id}
                movie={{
                  movie_id: card.id,
                  title: card.title,
                  poster_url: card.poster_url,
                }}
              />
            ) : (
              <BookCard
                key={card.id}
                book={{
                  book_id: card.id,
                  title: card.title,
                  image_url: card.image_url,
                }}
              />
            )
          )}
        </div>
      )}

      {!isLoading && !isError && cards.length === 0 && (
        <p className="text-center py-20 text-muted-foreground">
          Add some items to your library to get recommendations.
        </p>
      )}
    </div>
  );
}