'use client';

import { useState, useEffect } from 'react';
import api from '@/lib/api';
import MovieCard from '@/components/movie-card';
import BookCard from '@/components/book-card';

type Tab = 'movies' | 'books';

export default function HomePage() {
  const [tab, setTab] = useState<Tab>('movies');
  const [loading, setLoading] = useState(false);
  const [cards, setCards] = useState<any[]>([]);
  const [error, setError] = useState('');

  const fetchRecommendations = async () => {
    setLoading(true);
    setError('');
    setCards([]);
    try {
      // 1. Создаём задачу
      const direction = tab === 'movies' ? 'book_to_movie' : 'movie_to_book';
      const { data: task } = await api.post('/api/recommend', {
        selected_ids: [],
        direction,
        weights: { genre: 0.3, text: 0.4, image: 0.3 },
      });
      const taskId = task.task_id;

      // 2. Опросим результат
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
        setError('No recommendations yet.');
        setLoading(false);
        return;
      }

      // 3. Загрузим метаданные
      const ids = result.movies;
      const type = tab === 'movies' ? 'movie' : 'book';
      const { data: batch } = await api.get(`/api/items/batch?ids=${ids.join(',')}&type=${type}`);
      const items = batch.items || [];
      // Приведём к нужным полям
      const processed = items.map((item: any) => ({
        id: item.id,
        title: item.title,
        poster_url: item.image || '',
        image_url: item.image || '',
      }));
      setCards(processed);
    } catch (err: any) {
      setError(err.message || 'Failed');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchRecommendations();
  }, [tab]);

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

      {loading && <div className="flex justify-center py-20">Loading...</div>}
      {error && <p className="text-center py-20 text-destructive">{error}</p>}

      {cards.length > 0 && (
        <div className="grid grid-cols-2 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-2 p-4">
          {cards.map((card: any) => (
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
          ))}
        </div>
      )}

      {!loading && !error && cards.length === 0 && (
        <p className="text-center py-20 text-muted-foreground">
          Add some items to your library to get recommendations.
        </p>
      )}
    </div>
  );
}