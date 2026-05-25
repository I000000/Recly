'use client';

import React, { useState, useMemo } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { useRouter } from 'next/navigation';
import Link from 'next/link';
import { Trash2, BookOpen, Film, History, Bookmark, ChevronRight, ChevronLeft, Loader2 } from 'lucide-react';
import api from '@/lib/api';
import MovieCard from '@/components/movie-card';
import BookCard from '@/components/book-card';
import { Button } from '@/components/ui/button';

interface HistoryEntry {
  id: string;
  task_id: string;
  direction: string;
  selected_ids: string | string[];
  result: string;
  created_at: string;
}

interface SavedRecommendation {
  id: string;
  from_type: 'book' | 'movie';
  from_id: string;
  to_type: 'book' | 'movie';
  to_id: string;
}

export default function ProfilePage() {
  const router = useRouter();
  const queryClient = useQueryClient();

  const [showAllSaved, setShowAllSaved] = useState(false);
  // showAllHistory больше не нужен

  const { data: likedBooksRaw } = useQuery<string[]>({
    queryKey: ['likedBooks'],
    queryFn: async () => (await api.get('/api/user/library/books')).data.books.map((b: any) => b.book_id),
    staleTime: 1000 * 60 * 30,
  });
  const { data: likedMoviesRaw } = useQuery<string[]>({
    queryKey: ['likedMovies'],
    queryFn: async () => (await api.get('/api/user/library/movies')).data.movies.map((m: any) => m.movie_id),
    staleTime: 1000 * 60 * 30,
  });

  const {
    data: historyRaw,
    isLoading: historyLoading,
    isError: historyError,
  } = useQuery<HistoryEntry[]>({
    queryKey: ['history'],
    queryFn: async () => (await api.get('/api/user/recommendations/history')).data.history,
  });

  const {
    data: savedRaw,
    isLoading: savedLoading,
    isError: savedError,
  } = useQuery<SavedRecommendation[]>({
    queryKey: ['saved'],
    queryFn: async () => (await api.get('/api/user/recommendations/saved')).data.saved,
  });

  const likedBooks = likedBooksRaw ?? [];
  const likedMovies = likedMoviesRaw ?? [];
  const history = historyRaw ?? [];
  const saved = savedRaw ?? [];

  // Сбор уникальных рекомендованных ID
  const recommendedIds = useMemo(() => {
    const idsSet = new Set<string>();
    history.forEach(entry => {
      try {
        const results = JSON.parse(entry.result || '[]');
        (results as string[]).forEach(id => idsSet.add(id));
      } catch {}
    });
    return Array.from(idsSet).slice(0, 40);
  }, [history]);

  // Метаданные для сохранённых
  const savedIds = new Set<string>();
  saved.forEach(item => {
    savedIds.add(item.from_id);
    savedIds.add(item.to_id);
  });
  const {
    data: savedMeta = {},
    isLoading: savedMetaLoading,
    isError: savedMetaError,
  } = useQuery<Record<string, any>>({
    queryKey: ['savedMeta', [...savedIds].sort()],
    queryFn: async () => {
      if (savedIds.size === 0) return {};
      const res = await api.get(`/api/items/batch?ids=${[...savedIds].join(',')}&type=all`);
      const map: Record<string, any> = {};
      (res.data.items || []).forEach((item: any) => (map[item.id] = item));
      return map;
    },
    enabled: savedIds.size > 0,
  });

  // Метаданные для рекомендованных ID
  const {
    data: historyMeta = {},
    isLoading: historyMetaLoading,
    isError: historyMetaError,
  } = useQuery<Record<string, any>>({
    queryKey: ['historyMeta', recommendedIds],
    queryFn: async () => {
      if (recommendedIds.length === 0) return {};
      const res = await api.get(`/api/items/batch?ids=${recommendedIds.join(',')}&type=all`);
      const map: Record<string, any> = {};
      (res.data.items || []).forEach((item: any) => (map[item.id] = item));
      return map;
    },
    enabled: recommendedIds.length > 0,
  });

  const deleteSaved = useMutation({
    mutationFn: (id: string) => api.delete(`/api/user/recommendations/saved/${id}`),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['saved'] }),
  });

  const handleLogout = () => {
    localStorage.removeItem('token');
    document.cookie = 'token=; path=/; expires=Thu, 01 Jan 1970 00:00:00 UTC';
    router.push('/login');
  };

  const renderRecommendationCard = (toItem: any, onDelete?: () => void, cardKey?: string) => {
    if (!toItem) return null;
    return (
      <div key={cardKey} className="relative flex-shrink-0 w-36 group">
        {toItem.type === 'movie' ? (
          <MovieCard movie={{ movie_id: toItem.id, title: toItem.title, poster_url: toItem.image }} aspectRatio="2/3" />
        ) : (
          <BookCard book={{ book_id: toItem.id, title: toItem.title, image_url: toItem.image }} aspectRatio="2/3" />
        )}
        {onDelete && (
          <button
            onClick={(e) => {
              e.preventDefault();
              onDelete();
            }}
            className="absolute top-1 right-1 p-1 rounded-full bg-background/70 hover:bg-destructive/20 text-muted-foreground hover:text-destructive transition"
            title="Remove"
          >
            <Trash2 className="w-3.5 h-3.5" />
          </button>
        )}
      </div>
    );
  };

  // Упрощённая секция только для истории (горизонтальный скролл + "See More")
  const renderHistorySection = (
    items: string[],
    renderItem: (id: string) => React.ReactNode,
    emptyMessage: string,
    isLoading: boolean,
    isError: boolean
  ) => {
    if (isLoading) {
      return (
        <div className="flex justify-center py-12">
          <Loader2 className="w-6 h-6 animate-spin text-muted-foreground" />
        </div>
      );
    }
    if (isError) {
      return <p className="text-destructive text-center py-4">Failed to load data</p>;
    }
    if (items.length === 0) {
      return <p className="text-muted-foreground py-4">{emptyMessage}</p>;
    }

    const displayedItems = items.slice(0, 10);
    return (
      <div className="flex gap-2 overflow-x-auto pb-2 no-scrollbar">
        {displayedItems.map(id => renderItem(id))}
        <Link
          href="/history"
          className="flex-shrink-0 w-36 relative rounded-xl overflow-hidden border shadow-md group cursor-pointer"
        >
          <div className="w-full aspect-[2/3] bg-black/60 flex flex-col items-center justify-center gap-1 text-white">
            <BookOpen className="w-6 h-6" />
            <span className="text-xs font-semibold">See More</span>
          </div>
        </Link>
      </div>
    );
  };

  const renderScrollableSection = (
    items: any[],
    renderItem: (item: any, index: number) => React.ReactNode,
    showAll: boolean,
    toggleShowAll: () => void,
    emptyMessage: string,
    isLoading: boolean,
    isError: boolean
  ) => {
    if (isLoading) {
      return (
        <div className="flex justify-center py-12">
          <Loader2 className="w-6 h-6 animate-spin text-muted-foreground" />
        </div>
      );
    }
    if (isError) {
      return <p className="text-destructive text-center py-4">Failed to load data</p>;
    }
    if (items.length === 0) {
      return <p className="text-muted-foreground py-4">{emptyMessage}</p>;
    }

    const displayedItems = showAll ? items : items.slice(0, 10);
    return (
      <div>
        {!showAll ? (
          <div className="flex gap-2 overflow-x-auto pb-2 no-scrollbar">
            {displayedItems.map((item, idx) => renderItem(item, idx))}
          </div>
        ) : (
          <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 gap-2">
            {displayedItems.map((item, idx) => renderItem(item, idx))}
          </div>
        )}
        {items.length > 10 && (
          <button onClick={toggleShowAll} className="mt-3 text-sm text-primary hover:underline flex items-center gap-1">
            {showAll ? 'Show less' : 'See all'}
            {showAll ? <ChevronLeft className="w-4 h-4" /> : <ChevronRight className="w-4 h-4" />}
          </button>
        )}
      </div>
    );
  };

  return (
    <div className="min-h-screen pb-20">
      <div className="px-4 pt-6 pb-2">
        <h1 className="text-2xl font-bold tracking-tight">Your Profile</h1>
      </div>

      <div className="px-4 py-4">
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
          <div className="bg-secondary/30 rounded-xl p-4 flex items-center gap-3">
            <Film className="w-5 h-5 text-primary" />
            <div>
              <p className="text-2xl font-bold">{likedMovies.length}</p>
              <p className="text-xs text-muted-foreground">Movies</p>
            </div>
          </div>
          <div className="bg-secondary/30 rounded-xl p-4 flex items-center gap-3">
            <BookOpen className="w-5 h-5 text-primary" />
            <div>
              <p className="text-2xl font-bold">{likedBooks.length}</p>
              <p className="text-xs text-muted-foreground">Books</p>
            </div>
          </div>
          <div className="bg-secondary/30 rounded-xl p-4 flex items-center gap-3">
            <History className="w-5 h-5 text-primary" />
            <div>
              <p className="text-2xl font-bold">{history.length}</p>
              <p className="text-xs text-muted-foreground">Requests</p>
            </div>
          </div>
          <div className="bg-secondary/30 rounded-xl p-4 flex items-center gap-3">
            <Bookmark className="w-5 h-5 text-primary" />
            <div>
              <p className="text-2xl font-bold">{saved.length}</p>
              <p className="text-xs text-muted-foreground">Saved</p>
            </div>
          </div>
        </div>
      </div>

      {/* Сохранённые рекомендации */}
      <div className="px-4 py-4">
        <h2 className="text-xl font-semibold mb-4">Saved Recommendations</h2>
        {renderScrollableSection(
          saved,
          (item: SavedRecommendation) => {
            const toItem = savedMeta[item.to_id];
            return renderRecommendationCard(toItem, () => deleteSaved.mutate(item.id), item.id);
          },
          showAllSaved,
          () => setShowAllSaved(!showAllSaved),
          'No saved recommendations yet.',
          savedLoading || savedMetaLoading,
          savedError || savedMetaError
        )}
      </div>

      {/* История (упрощённая) */}
      <div className="px-4 py-4">
        <h2 className="text-xl font-semibold mb-4">Recent Activity</h2>
        {renderHistorySection(
          recommendedIds,
          (id: string) => {
            const item = historyMeta[id];
            return renderRecommendationCard(item, undefined, id);
          },
          'No recommendations yet.',
          historyLoading || historyMetaLoading,
          historyError || historyMetaError
        )}
      </div>

      <div className="px-4 pb-8">
        <Button onClick={handleLogout} variant="destructive" className="w-full">
          Log out
        </Button>
      </div>
    </div>
  );
}