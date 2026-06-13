'use client';

import { useQuery } from '@tanstack/react-query';
import { useRouter } from 'next/navigation';
import { ArrowLeft, Loader2 } from 'lucide-react';
import api from '@/lib/api';
import MovieCard from '@/components/movie-card';
import BookCard from '@/components/book-card';

interface ViewedItem {
  id: string;
  item_id: string;
  item_type: 'book' | 'movie';
  viewed_at: string;
}

export default function ViewedPage() {
  const router = useRouter();

  const { data: views, isLoading } = useQuery<ViewedItem[]>({
    queryKey: ['views'],
    queryFn: async () => (await api.get('/api/user/views')).data.views || [],
    staleTime: 2 * 60 * 1000,
  });

  const itemIds = views?.map(v => v.item_id) || [];
  const { data: itemsMeta = {}, isLoading: metaLoading } = useQuery({
    queryKey: ['viewedMeta', itemIds],
    queryFn: async () => {
      if (itemIds.length === 0) return {};
      const res = await api.get(`/api/items/batch?ids=${itemIds.join(',')}&type=all`);
      const map: Record<string, any> = {};
      (res.data.items || []).forEach((item: any) => (map[item.id] = item));
      return map;
    },
    enabled: itemIds.length > 0,
  });

  if (isLoading || metaLoading) {
    return (
      <div className="min-h-screen flex justify-center items-center">
        <Loader2 className="w-8 h-8 animate-spin" />
      </div>
    );
  }

  const sorted = views ? [...views].sort((a, b) => new Date(b.viewed_at).getTime() - new Date(a.viewed_at).getTime()) : [];

  return (
    <div className="min-h-screen pb-20">
      <div className="px-4 pt-6 pb-2 flex items-center gap-3">
        <button onClick={() => router.back()} className="p-2 rounded-full hover:bg-secondary">
          <ArrowLeft className="w-5 h-5" />
        </button>
        <h1 className="text-2xl font-bold tracking-tight">Recently viewed</h1>
      </div>
      {sorted.length === 0 ? (
        <p className="text-center text-muted-foreground py-20">No viewed items yet.</p>
      ) : (
        <div className="grid grid-cols-2 sm:grid-cols-[repeat(auto-fill,minmax(225px,1fr))] gap-2 p-4">
          {sorted.map(view => {
            const item = itemsMeta[view.item_id];
            if (!item) return null;
            return view.item_type === 'movie' ? (
              <MovieCard key={view.id} movie={{ movie_id: item.id, title: item.title, poster_url: item.image }} aspectRatio="2/3" />
            ) : (
              <BookCard key={view.id} book={{ book_id: item.id, title: item.title, image_url: item.image }} aspectRatio="2/3" />
            );
          })}
        </div>
      )}
    </div>
  );
}