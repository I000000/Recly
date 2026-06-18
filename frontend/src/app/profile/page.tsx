'use client';

import { useMemo, useRef, useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { useRouter } from 'next/navigation';
import Link from 'next/link';
import { BookOpen, Film, History, Bookmark, Loader2, Settings, Camera } from 'lucide-react';
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

interface ViewedItem {
  id: string;
  item_id: string;
  item_type: 'book' | 'movie';
  viewed_at: string;
}

interface SavedItem {
  id: string;
  item_id: string;
  item_type: 'book' | 'movie';
  saved_at: string;
}

export default function ProfilePage() {
  const router = useRouter();
  const queryClient = useQueryClient();
  const fileInputRef = useRef<HTMLInputElement>(null);
  const [uploading, setUploading] = useState(false);

  const { data: profile, isLoading: profileLoading } = useQuery({
    queryKey: ['userProfile'],
    queryFn: async () => (await api.get('/api/user/profile')).data,
  });

  const { data: likedMoviesRaw } = useQuery<string[]>({
    queryKey: ['likedMovies'],
    queryFn: async () => (await api.get('/api/user/library/movies')).data.movies.map((m: any) => m.movie_id),
    staleTime: 1000 * 60 * 30,
  });
  const { data: likedBooksRaw } = useQuery<string[]>({
    queryKey: ['likedBooks'],
    queryFn: async () => (await api.get('/api/user/library/books')).data.books.map((b: any) => b.book_id),
    staleTime: 1000 * 60 * 30,
  });
  const likedMovies = likedMoviesRaw ?? [];
  const likedBooks = likedBooksRaw ?? [];

  const { data: savedRaw, isLoading: savedLoading, isError: savedError } = useQuery<SavedItem[]>({
    queryKey: ['saved'],
    queryFn: async () => (await api.get('/api/user/saved-items')).data.saved,
    staleTime: 5 * 60 * 1000,
  });
  const saved = savedRaw ?? [];

  const { data: historyRaw, isLoading: historyLoading, isError: historyError } = useQuery<HistoryEntry[]>({
    queryKey: ['history'],
    queryFn: async () => (await api.get('/api/user/recommendations/history')).data.history,
    staleTime: 5 * 60 * 1000,
  });
  const history = historyRaw ?? [];

  const { data: viewedRaw, isLoading: viewedLoading, isError: viewedError } = useQuery<ViewedItem[]>({
    queryKey: ['views'],
    queryFn: async () => (await api.get('/api/user/views')).data.views || [],
    staleTime: 2 * 60 * 1000,
  });

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

  const viewedIds = useMemo(() => {
    if (!viewedRaw) return [];
    const sorted = [...viewedRaw].sort((a, b) => new Date(b.viewed_at).getTime() - new Date(a.viewed_at).getTime());
    return sorted.map(v => v.item_id);
  }, [viewedRaw]);

  const savedIds = saved.map(item => item.item_id);

  const { data: savedMeta = {}, isLoading: savedMetaLoading, isError: savedMetaError } = useQuery<Record<string, any>>({
    queryKey: ['savedMeta', savedIds],
    queryFn: async () => {
      if (savedIds.length === 0) return {};
      const res = await api.get(`/api/items/batch?ids=${savedIds.join(',')}&type=all`);
      const map: Record<string, any> = {};
      (res.data.items || []).forEach((item: any) => (map[item.id] = item));
      return map;
    },
    enabled: savedIds.length > 0,
    staleTime: 10 * 60 * 1000,
  });

  const { data: historyMeta = {}, isLoading: historyMetaLoading, isError: historyMetaError } = useQuery<Record<string, any>>({
    queryKey: ['historyMeta', recommendedIds],
    queryFn: async () => {
      if (recommendedIds.length === 0) return {};
      const res = await api.get(`/api/items/batch?ids=${recommendedIds.join(',')}&type=all`);
      const map: Record<string, any> = {};
      (res.data.items || []).forEach((item: any) => (map[item.id] = item));
      return map;
    },
    enabled: recommendedIds.length > 0,
    staleTime: 10 * 60 * 1000,
  });

  const { data: viewedMeta = {}, isLoading: viewedMetaLoading, isError: viewedMetaError } = useQuery<Record<string, any>>({
    queryKey: ['viewedMeta', viewedIds],
    queryFn: async () => {
      if (viewedIds.length === 0) return {};
      const res = await api.get(`/api/items/batch?ids=${viewedIds.join(',')}&type=all`);
      const map: Record<string, any> = {};
      (res.data.items || []).forEach((item: any) => (map[item.id] = item));
      return map;
    },
    enabled: viewedIds.length > 0,
    staleTime: 10 * 60 * 1000,
  });

  const deleteSaved = useMutation({
    mutationFn: (id: string) => api.delete(`/api/user/saved-items/${id}`),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['saved'] }),
  });

  const uploadAvatar = useMutation({
    mutationFn: async (file: File) => {
      const formData = new FormData();
      formData.append('avatar', file);
      const res = await api.post('/api/user/avatar', formData, {
        headers: { 'Content-Type': 'multipart/form-data' },
      });
      return res.data.avatar_url;
    },
    onSuccess: (newAvatarUrl) => {
      queryClient.setQueryData(['userProfile'], (old: any) => ({ ...old, avatar_url: newAvatarUrl }));
      queryClient.invalidateQueries({ queryKey: ['userProfile'] });
      setUploading(false);
    },
    onError: () => setUploading(false),
  });

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file) {
      setUploading(true);
      uploadAvatar.mutate(file);
    }
    if (fileInputRef.current) fileInputRef.current.value = '';
  };

  const renderCard = (item: any) => {
    if (!item) return null;
    return (
      <div className="relative flex-shrink-0 w-36 group">
        {item.type === 'movie' ? (
          <MovieCard movie={{ movie_id: item.id, title: item.title, poster_url: item.image }} aspectRatio="2/3" />
        ) : (
          <BookCard book={{ book_id: item.id, title: item.title, image_url: item.image }} aspectRatio="2/3" />
        )}
      </div>
    );
  };

  const renderHorizontalSection = (
    itemIds: string[],
    itemsMeta: Record<string, any>,
    emptyMessage: string,
    isLoading: boolean,
    isError: boolean,
    seeMorePath: string,
  ) => {
    if (isLoading) {
      return <div className="flex justify-center py-12"><Loader2 className="w-6 h-6 animate-spin text-muted-foreground" /></div>;
    }
    if (isError) {
      return <p className="text-destructive text-center py-4">Failed to load data</p>;
    }
    if (itemIds.length === 0) {
      return <p className="text-muted-foreground py-4">{emptyMessage}</p>;
    }

    const hasMore = itemIds.length > 6;
    const visibleIds = hasMore ? itemIds.slice(0, 6) : itemIds;
    const nextId = hasMore ? itemIds[6] : null;
    const nextItem = nextId ? itemsMeta[nextId] : null;

    return (
      <div className="scroll-container pb-2">
        <div className="flex gap-2 min-w-max">
          {visibleIds.map(id => <div key={id} className="w-36 flex-shrink-0">{renderCard(itemsMeta[id])}</div>)}
          {hasMore && nextItem && (
            <div
              onClick={() => router.push(seeMorePath)}
              className="flex-shrink-0 w-36 relative rounded-xl overflow-hidden border shadow-md group cursor-pointer"
            >
              <div className="w-full aspect-[2/3] relative">
                {nextItem.image ? (
                  <img src={nextItem.image} alt={nextItem.title || ''} className="absolute inset-0 w-full h-full object-cover" />
                ) : (
                  <div className="absolute inset-0 bg-muted flex items-center justify-center text-xs">No poster</div>
                )}
                <div className="absolute bottom-0 left-0 right-0 bg-gradient-to-t from-black/70 via-black/20 to-transparent p-3 pt-22">
                  <h3 className="text-white text-sm font-semibold leading-snug line-clamp-2 drop-shadow-md">
                    {nextItem.title || 'Untitled'}
                  </h3>
                </div>
              </div>
              <div className="absolute inset-0 bg-black/60 flex flex-col items-center justify-center gap-1 text-white opacity-100 group-hover:opacity-90 transition-opacity">
                <span className="text-xs font-semibold">See More</span>
              </div>
            </div>
          )}
        </div>
      </div>
    );
  };

  if (profileLoading) {
    return <div className="min-h-screen flex justify-center items-center"><Loader2 className="w-8 h-8 animate-spin" /></div>;
  }

  return (
    <div className="min-h-screen pb-20 mx-auto max-w-full overflow-x-hidden">
      <div className="px-4 pt-6 pb-2 flex justify-between items-center">
        <h1 className="text-2xl font-bold tracking-tight">Your Profile</h1>
        <Link href="/settings" className="p-2 rounded-full hover:bg-secondary md:hidden">
          <Settings className="w-5 h-5" />
        </Link>
      </div>

      <div className="px-4 py-4 flex flex-col md:flex-row md:items-center gap-6">
        <div className="flex flex-col items-center gap-2">
          <div onClick={() => fileInputRef.current?.click()} className="relative w-28 h-28 rounded-full overflow-hidden bg-secondary cursor-pointer">
            {profile?.avatar_url ? (
              <img src={profile.avatar_url} alt="Avatar" className="w-full h-full object-cover" />
            ) : (
              <div className="w-full h-full flex items-center justify-center text-5xl font-bold text-muted-foreground">
                {profile?.name?.[0]?.toUpperCase() || 'U'}
              </div>
            )}
            <input ref={fileInputRef} type="file" accept="image/jpeg,image/png,image/webp" onChange={handleFileChange} className="hidden" />
          </div>
          <div className="text-center"><p className="font-semibold text-xl">{profile?.name || 'User'}</p></div>
        </div>

        <div className="grid grid-cols-2 gap-4">
          <Link href="/library?tab=movies" className="bg-secondary/30 rounded-xl p-4 hover:bg-secondary/50 transition flex items-center gap-3">
            <Film className="w-7 h-7 text-primary" />
            <div><div className="text-2xl font-bold">{likedMovies.length}</div><div className="text-sm text-muted-foreground">Movies</div></div>
          </Link>
          <Link href="/library?tab=books" className="bg-secondary/30 rounded-xl p-4 hover:bg-secondary/50 transition flex items-center gap-3">
            <BookOpen className="w-7 h-7 text-primary" />
            <div><div className="text-2xl font-bold">{likedBooks.length}</div><div className="text-sm text-muted-foreground">Books</div></div>
          </Link>
          <Link href="/viewed" className="bg-secondary/30 rounded-xl p-4 hover:bg-secondary/50 transition flex items-center gap-3">
            <History className="w-7 h-7 text-primary" />
            <div><div className="text-2xl font-bold">{viewedIds.length}</div><div className="text-sm text-muted-foreground">Viewed</div></div>
          </Link>
          <Link href="/saved" className="bg-secondary/30 rounded-xl p-4 hover:bg-secondary/50 transition flex items-center gap-3">
            <Bookmark className="w-7 h-7 text-primary" />
            <div><div className="text-2xl font-bold">{saved.length}</div><div className="text-sm text-muted-foreground">Saved</div></div>
          </Link>
        </div>
      </div>

      {(savedIds.length > 0 || savedLoading || savedMetaLoading || savedError || savedMetaError) && (
        <div className="px-4 py-4">
          <h2 className="text-xl font-semibold mb-4">Saved</h2>
          {renderHorizontalSection(
            savedIds,
            savedMeta,
            'No saved items yet.',
            savedLoading || savedMetaLoading,
            savedError || savedMetaError,
            '/saved',
          )}
        </div>
      )}
      {(viewedIds.length > 0 || viewedLoading || viewedMetaLoading || viewedError || viewedMetaError) && (
        <div className="px-4 py-4">
          <h2 className="text-xl font-semibold mb-4">Recently Viewed</h2>
          {renderHorizontalSection(
            viewedIds,
            viewedMeta,
            'No views yet.',
            viewedLoading || viewedMetaLoading,
            viewedError || viewedMetaError,
            '/viewed',
          )}
        </div>
      )}
      {(recommendedIds.length > 0 || historyLoading || historyMetaLoading || historyError || historyMetaError) && (
        <div className="px-4 py-4">
          <h2 className="text-xl font-semibold mb-4">Recently Recommended</h2>
          {renderHorizontalSection(
            recommendedIds,
            historyMeta,
            'No recommendations yet.',
            historyLoading || historyMetaLoading,
            historyError || historyMetaError,
            '/history',
          )}
        </div>
      )}
    </div>
  );
}