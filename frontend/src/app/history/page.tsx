'use client';

import { useMemo } from 'react';
import { useQuery } from '@tanstack/react-query';
import Link from 'next/link';
import { ArrowLeft, Loader2 } from 'lucide-react';
import api from '@/lib/api';
import MovieCard from '@/components/movie-card';
import BookCard from '@/components/book-card';

interface HistoryEntry {
  id: string;
  task_id: string;
  direction: string;
  selected_ids: string | string[];
  result: string;
  created_at: string;
}

export default function FullHistoryPage() {
  const { data: historyRaw, isLoading, isError } = useQuery<HistoryEntry[]>({
    queryKey: ['history'],
    queryFn: async () => (await api.get('/api/user/recommendations/history')).data.history,
  });

  const history = historyRaw ?? [];

  const recommendedIds = useMemo(() => {
    const idsSet = new Set<string>();
    history.forEach(entry => {
      try {
        const results = JSON.parse(entry.result || '[]');
        (results as string[]).forEach(id => idsSet.add(id));
      } catch {}
    });
    return Array.from(idsSet);
  }, [history]);

  const { data: itemsMap = {}, isLoading: metaLoading } = useQuery<Record<string, any>>({
    queryKey: ['historyMetaFull', recommendedIds],
    queryFn: async () => {
      if (recommendedIds.length === 0) return {};
      const res = await api.get(`/api/items/batch?ids=${recommendedIds.join(',')}&type=all`);
      const map: Record<string, any> = {};
      (res.data.items || []).forEach((item: any) => (map[item.id] = item));
      return map;
    },
    enabled: recommendedIds.length > 0,
  });

  if (isLoading || metaLoading) {
    return (
      <div className="min-h-screen flex justify-center items-center">
        <Loader2 className="w-8 h-8 animate-spin" />
      </div>
    );
  }

  if (isError) {
    return <div className="p-8 text-destructive">Failed to load history.</div>;
  }

  return (
    <div className="min-h-screen pb-20">
      <div className="px-4 pt-6 pb-2 flex items-center gap-3">
        <Link href="/profile" className="p-2 rounded-full hover:bg-secondary">
          <ArrowLeft className="w-5 h-5" />
        </Link>
        <h1 className="text-2xl font-bold tracking-tight">All Recommendations</h1>
      </div>

      {recommendedIds.length === 0 ? (
        <p className="text-muted-foreground text-center py-20">No recommendations yet.</p>
      ) : (
        <div className="grid grid-cols-2 sm:grid-cols-[repeat(auto-fill,minmax(225px,1fr))] gap-2 p-4">
          {recommendedIds.map(id => {
            const item = itemsMap[id];
            if (!item) return null;
            return item.type === 'movie' ? (
              <MovieCard
                key={id}
                movie={{ movie_id: id, title: item.title, poster_url: item.image }}
                aspectRatio="2/3"
              />
            ) : (
              <BookCard
                key={id}
                book={{ book_id: id, title: item.title, image_url: item.image }}
                aspectRatio="2/3"
              />
            );
          })}
        </div>
      )}
    </div>
  );
}