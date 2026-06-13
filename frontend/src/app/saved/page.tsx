'use client';

import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { useRouter } from 'next/navigation';
import { ArrowLeft, Loader2, Trash2 } from 'lucide-react';
import api from '@/lib/api';
import MovieCard from '@/components/movie-card';
import BookCard from '@/components/book-card';

export default function SavedPage() {
  const router = useRouter();
  const queryClient = useQueryClient();

  const { data: savedRaw, isLoading, isError } = useQuery<any[]>({
    queryKey: ['saved'],
    queryFn: async () => (await api.get('/api/user/saved-items')).data.saved,
    staleTime: 5 * 60 * 1000,
  });
  const saved = savedRaw ?? [];

  const savedIds = saved.map(item => item.item_id);
  const { data: savedMeta = {}, isLoading: metaLoading } = useQuery<Record<string, any>>({
    queryKey: ['savedMeta', savedIds],
    queryFn: async () => {
      if (savedIds.length === 0) return {};
      const res = await api.get(`/api/items/batch?ids=${savedIds.join(',')}&type=all`);
      const map: Record<string, any> = {};
      (res.data.items || []).forEach((item: any) => (map[item.id] = item));
      return map;
    },
    enabled: savedIds.length > 0,
  });

  const deleteSaved = useMutation({
    mutationFn: (id: string) => api.delete(`/api/user/saved-items/${id}`),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['saved'] });
    },
  });

  if (isLoading || metaLoading) {
    return (
      <div className="min-h-screen flex justify-center items-center">
        <Loader2 className="w-8 h-8 animate-spin" />
      </div>
    );
  }
  if (isError) {
    return <div className="p-8 text-destructive text-center">Failed to load saved items.</div>;
  }

  return (
    <div className="min-h-screen pb-20">
      <div className="px-4 pt-6 pb-2 flex items-center gap-3">
        <button onClick={() => router.back()} className="p-2 rounded-full hover:bg-secondary">
          <ArrowLeft className="w-5 h-5" />
        </button>
        <h1 className="text-2xl font-bold tracking-tight">Saved items</h1>
      </div>
      {saved.length === 0 ? (
        <p className="text-center text-muted-foreground py-20">No saved items yet.</p>
      ) : (
        <div className="grid grid-cols-2 sm:grid-cols-[repeat(auto-fill,minmax(225px,1fr))] gap-2 p-4">
          {saved.map(item => {
            const toItem = savedMeta[item.item_id];
            if (!toItem) return null;
            return (
              <div key={item.id} className="relative group">
                {toItem.type === 'movie' ? (
                  <MovieCard movie={{ movie_id: toItem.id, title: toItem.title, poster_url: toItem.image }} aspectRatio="2/3" />
                ) : (
                  <BookCard book={{ book_id: toItem.id, title: toItem.title, image_url: toItem.image }} aspectRatio="2/3" />
                )}
                <button
                  onClick={() => deleteSaved.mutate(item.id)}
                  className="absolute top-2 right-2 p-1.5 rounded-full bg-destructive/80 text-white hover:bg-destructive transition"
                >
                  <Trash2 className="w-4 h-4" />
                </button>
              </div>
            );
          })}
        </div>
      )}
    </div>
  );
}