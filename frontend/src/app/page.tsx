'use client';

import { useState, useEffect, useRef, useCallback } from 'react';
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { Loader2 } from 'lucide-react';
import api from '@/lib/api';
import { loadMovies, loadBooks } from '@/lib/data';
import { RecommendationResult, Movie, Book } from '@/types';
import MovieCard from '@/components/movie-card';
import BookCard from '@/components/book-card';

type Tab = 'movies' | 'books';

// Вспомогательный тип для различения состояния
interface DoneResult extends RecommendationResult {
  status: 'done';
  movies: string[];
}

export default function HomePage() {
  const [tab, setTab] = useState<Tab>('movies');
  const [taskId, setTaskId] = useState<string | null>(null);
  const [moviesMap, setMoviesMap] = useState<Record<string, Movie>>({});
  const [booksMap, setBooksMap] = useState<Record<string, Book>>({});
  const queryClient = useQueryClient();
  const hasRequested = useRef(false);

  // Pull‑to‑refresh
  const touchStartY = useRef(0);
  const [refreshing, setRefreshing] = useState(false);

  // Загружаем метаданные один раз
  useEffect(() => {
    loadMovies().then((movies) => {
      const map: Record<string, Movie> = {};
      movies.forEach((m: any) => { map[String(m.movie_id)] = m; });
      setMoviesMap(map);
    });
    loadBooks().then((books) => {
      const map: Record<string, Book> = {};
      books.forEach((b: any) => { map[String(b.book_id)] = b; });
      setBooksMap(map);
    });
  }, []);

  // Мутация для отправки задачи
  const mutation = useMutation({
    mutationFn: async (currentTab: Tab) => {
      const direction = currentTab === 'movies' ? 'movie_to_movie' : 'book_to_book';
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

  // Проверяем кэш под общим ключом (без taskId)
  useEffect(() => {
    const cached = queryClient.getQueryData<RecommendationResult>(['recommendation', tab]);
    hasRequested.current = false;
    if (!cached || cached.status !== 'done') {
      if (hasRequested.current) return;
      hasRequested.current = true;
      setTaskId(null);
      mutation.mutate(tab);
    } else {
      queryClient.setQueryData(['recommendation', tab, taskId], cached);
    }
  }, [tab]);

  // Опрос результата (ключ с taskId, чтобы триггерился при новом запросе)
  const { data: rawResult } = useQuery<RecommendationResult>({
    queryKey: ['recommendation', tab, taskId],
    queryFn: async () => {
      if (!taskId) {
        // Этот код никогда не выполнится из-за enabled, но нужен для типов
        return { status: 'pending' } as RecommendationResult;
      }
      const res = await api.get(`/api/result/${taskId}`);
      return res.data as RecommendationResult;
    },
    enabled: !!taskId,
    refetchInterval: (query) => {
      const data = query.state.data;
      if (data && data.status === 'pending') return 2000;
      return false;
    },
    staleTime: 1000 * 60 * 30,
    gcTime: 1000 * 60 * 30,
  });

  const result = rawResult as RecommendationResult | undefined;

  // При успешном получении данных переносим результат в общий кэш
  useEffect(() => {
    if (result && result.status === 'done') {
      queryClient.setQueryData(['recommendation', tab], result);
    }
  }, [result, tab, queryClient]);

  // Pull‑to‑refresh handler – обновляет только текущую вкладку
  const handleTouchStart = useCallback((e: React.TouchEvent) => {
    touchStartY.current = e.touches[0].clientY;
  }, []);

  const handleTouchEnd = useCallback((e: React.TouchEvent) => {
    const delta = e.changedTouches[0].clientY - touchStartY.current;
    if (delta > 80 && window.scrollY === 0) {
      setRefreshing(true);
      // Инвалидируем оба ключа
      queryClient.invalidateQueries({ queryKey: ['recommendation', tab] });
      queryClient.invalidateQueries({ queryKey: ['recommendation', tab, taskId] });
      setTaskId(null);
      hasRequested.current = false;
      mutation.mutate(tab);
      setTimeout(() => setRefreshing(false), 2000);
    }
  }, [tab, taskId, queryClient, mutation]);

  const isDone = result?.status === 'done';
  const displayResult = isDone ? (result as DoneResult) : null;

  return (
    <div
      className="min-h-screen pb-20"
      onTouchStart={handleTouchStart}
      onTouchEnd={handleTouchEnd}
    >
      {/* Логотип (не sticky) */}
      <div className="px-4 pt-6 pb-2">
        <h1 className="text-2xl font-bold tracking-tight">Recly</h1>
      </div>

      {/* Sticky панель с кнопками */}
      <div className="sticky top-0 z-10 bg-background/80 backdrop-blur-sm border-b border-border px-4 py-2">
        <div className="flex gap-2">
          <button
            onClick={() => setTab('movies')}
            className={`px-4 py-2 rounded-full text-sm font-medium transition-colors ${
              tab === 'movies'
                ? 'bg-primary text-primary-foreground'
                : 'bg-secondary text-secondary-foreground hover:bg-secondary/80'
            }`}
          >
            Movies
          </button>
          <button
            onClick={() => setTab('books')}
            className={`px-4 py-2 rounded-full text-sm font-medium transition-colors ${
              tab === 'books'
                ? 'bg-primary text-primary-foreground'
                : 'bg-secondary text-secondary-foreground hover:bg-secondary/80'
            }`}
          >
            Books
          </button>
        </div>
        {refreshing && (
          <div className="flex justify-center mt-2">
            <Loader2 className="w-5 h-5 animate-spin text-muted-foreground" />
          </div>
        )}
      </div>

      {/* Состояния загрузки */}
      {mutation.isPending && !displayResult && (
        <div className="flex justify-center py-20">
          <Loader2 className="w-8 h-8 animate-spin text-muted-foreground" />
        </div>
      )}
      {result?.status === 'pending' && !displayResult && (
        <div className="flex justify-center py-20">
          <Loader2 className="w-8 h-8 animate-spin text-muted-foreground" />
        </div>
      )}
      {result?.status === 'error' && (
        <div className="flex justify-center py-20 text-destructive">Something went wrong</div>
      )}

      {/* Сетка с карточками */}
      {displayResult && (
        <div className="grid grid-cols-2 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-2 p-4">
          {tab === 'movies' &&
            displayResult.movies?.map((id) => {
              const movie = moviesMap[String(id)];
              return movie ? <MovieCard key={id} movie={movie} /> : null;
            })}
          {tab === 'books' &&
            displayResult.movies?.map((id) => {
              const bookData = booksMap[String(id)];
              return bookData ? <BookCard key={id} book={bookData} /> : null;
            })}
        </div>
      )}
    </div>
  );
}