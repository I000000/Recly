'use client';

import { useState, useEffect, useRef } from 'react';
import { useInfiniteQuery, useQuery, useQueryClient } from '@tanstack/react-query';
import { Loader2 } from 'lucide-react';
import api from '@/lib/api';
import MovieCard from '@/components/movie-card';
import BookCard from '@/components/book-card';

type Tab = 'movies' | 'books';

export default function HomePage() {
  const [tab, setTab] = useState<Tab>('movies');
  const queryClient = useQueryClient();
  const sentinelRef = useRef<HTMLDivElement>(null);

  const { data: likedBooks = [] } = useQuery<string[]>({
    queryKey: ['likedBooks'],
    queryFn: async () => (await api.get('/api/user/library/books')).data.books.map((b: any) => b.book_id),
    staleTime: 1000 * 60 * 30,
  });
  const { data: likedMovies = [] } = useQuery<string[]>({
    queryKey: ['likedMovies'],
    queryFn: async () => (await api.get('/api/user/library/movies')).data.movies.map((m: any) => m.movie_id),
    staleTime: 1000 * 60 * 30,
  });

  useEffect(() => {
    const saved = sessionStorage.getItem('homeTab');
    if (saved === 'movies' || saved === 'books') {
      setTab(saved);
    }
  }, []);

  const handleTabChange = (newTab: Tab) => {
    setTab(newTab);
    sessionStorage.setItem('homeTab', newTab);
  };

  useEffect(() => {
    if (typeof window === 'undefined') return;
    const renderCount = parseInt(sessionStorage.getItem('home_render_count') || '0', 10);
    if (renderCount === 0) {
      queryClient.removeQueries({ queryKey: ['homeRecommendations'] });
    }
    sessionStorage.setItem('home_render_count', String(renderCount + 1));
  }, [queryClient]);

  const fetchRecommendations = async ({ pageParam }: { pageParam: string[] }) => {
    const excludeIds = pageParam ?? [];
    const direction = tab === 'movies' ? 'book_to_movie' : 'movie_to_book';

    const allLiked = [
      ...likedBooks.map(id => `book_${id}`),
      ...likedMovies.map(id => `movie_${id}`)
    ];
    const shuffled = [...allLiked].sort(() => Math.random() - 0.5);
    const selected = shuffled.slice(0, Math.min(7, shuffled.length));

    const { data: task } = await api.post('/api/recommend', {
      selected_ids: selected,
      direction,
      weights: { genre: 0.3, text: 0.4, image: 0.3 },
      exclude_ids: excludeIds,
    });
    const taskId = task.task_id;

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
      return { cards: [], excludeIds };
    }

    queryClient.invalidateQueries({ queryKey: ['history'] });

    const ids = result.movies;
    const type = tab === 'movies' ? 'movie' : 'book';
    const { data: batch } = await api.get(`/api/items/batch?ids=${ids.join(',')}&type=${type}`);
    const items = batch.items || [];

    const newCards = items.map((item: any) => ({
      id: item.id,
      title: item.title,
      poster_url: item.image || '',
      image_url: item.image || '',
      type: item.type,
    }));

    const newExcludeIds = [...excludeIds, ...newCards.map((c: any) => c.id)];

    return { cards: newCards, excludeIds: newExcludeIds };
  };

  const {
    data,
    fetchNextPage,
    hasNextPage,
    isFetchingNextPage,
    isFetching,
    isError,
    error,
  } = useInfiniteQuery({
    queryKey: ['homeRecommendations', tab],
    queryFn: fetchRecommendations,
    initialPageParam: [] as string[],
    getNextPageParam: (lastPage) => {
      if (lastPage.cards.length === 0) return undefined;
      return lastPage.excludeIds;
    },
    staleTime: Infinity,
    refetchOnMount: false,
    refetchOnWindowFocus: false,
    enabled: typeof window !== 'undefined',
  });

  const cards = data?.pages.flatMap(page => page.cards) ?? [];

  useEffect(() => {
    if (!hasNextPage || isFetchingNextPage) return;
    const sentinel = sentinelRef.current;
    if (!sentinel) return;
    const observer = new IntersectionObserver(
      (entries) => {
        if (entries[0].isIntersecting) fetchNextPage();
      },
      { rootMargin: '200px' }
    );
    observer.observe(sentinel);
    return () => observer.disconnect();
  }, [hasNextPage, isFetchingNextPage, fetchNextPage]);

  const isFirstLoading = !data && isFetching;

  return (
    <div className="min-h-screen pb-20">
      <div className="px-4 pt-6 pb-2">
        <h1 className="text-2xl font-bold tracking-tight">Recly</h1>
      </div>
      <div className="sticky top-0 z-10 bg-background/80 backdrop-blur-sm border-b border-border px-4 py-2">
        <div className="flex gap-2">
          <button onClick={() => handleTabChange('movies')} className={`px-4 py-2 rounded-full text-sm font-medium ${tab === 'movies' ? 'bg-primary text-primary-foreground' : 'bg-secondary'}`}>
            Movies
          </button>
          <button onClick={() => handleTabChange('books')} className={`px-4 py-2 rounded-full text-sm font-medium ${tab === 'books' ? 'bg-primary text-primary-foreground' : 'bg-secondary'}`}>
            Books
          </button>
        </div>
      </div>

      {isFirstLoading && typeof window !== 'undefined' && (
        <div className="flex justify-center py-20">
          <Loader2 className="w-8 h-8 animate-spin text-muted-foreground" />
        </div>
      )}

      {isError && <p className="text-center py-20 text-destructive">{(error as any)?.message || 'Failed'}</p>}

      {!isFirstLoading && cards.length > 0 && (
        <div className="grid grid-cols-2 sm:grid-cols-[repeat(auto-fill,minmax(225px,1fr))] gap-2 p-4">
          {cards.map((card: any) =>
            card.type === 'movie' ? (
              <MovieCard key={card.id} movie={{ movie_id: card.id, title: card.title, poster_url: card.poster_url }} />
            ) : (
              <BookCard key={card.id} book={{ book_id: card.id, title: card.title, image_url: card.image_url }} />
            )
          )}
        </div>
      )}

      {!isFirstLoading && !isError && cards.length === 0 && typeof window !== 'undefined' && (
        <p className="text-center py-20 text-muted-foreground">
          Add some items to your library to get recommendations.
        </p>
      )}

      <div ref={sentinelRef} className="h-1" />

      {isFetchingNextPage && (
        <div className="flex justify-center py-4">
          <div className="w-6 h-6 border-2 border-primary border-t-transparent rounded-full animate-spin" />
        </div>
      )}
    </div>
  );
}