'use client';

import { useState, useEffect } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { Search, Loader2, Plus, X } from 'lucide-react';
import api from '@/lib/api';
import MovieCard from '@/components/movie-card';
import BookCard from '@/components/book-card';
import ItemSelector from '@/components/item-selector';

export default function LibraryPage() {
  const queryClient = useQueryClient();
  const [activeTab, setActiveTab] = useState<'movies' | 'books'>('movies');
  const [showAddModal, setShowAddModal] = useState(false);
  const [libraryQuery, setLibraryQuery] = useState('');

  // Восстанавливаем сохранённую вкладку после монтирования
  useEffect(() => {
    const saved = sessionStorage.getItem('libraryTab');
    if (saved === 'movies' || saved === 'books') {
      setActiveTab(saved);
    }
  }, []);

  // Функция переключения вкладки с сохранением
  const handleTabChange = (tab: 'movies' | 'books') => {
    setActiveTab(tab);
    sessionStorage.setItem('libraryTab', tab);
  };

  const { data: likedBooks = [] } = useQuery<string[]>({
    queryKey: ['likedBooks'],
    queryFn: async () => (await api.get('/api/user/library/books')).data.books.map((b: any) => b.book_id),
    staleTime: 1000 * 60 * 30,
    gcTime: 60 * 60 * 1000,
    enabled: activeTab === 'books',
  });

  const { data: likedMovies = [] } = useQuery<string[]>({
    queryKey: ['likedMovies'],
    queryFn: async () => (await api.get('/api/user/library/movies')).data.movies.map((m: any) => m.movie_id),
    staleTime: 1000 * 60 * 30,
    gcTime: 60 * 60 * 1000,
    enabled: activeTab === 'movies',
  });

  const ids = activeTab === 'movies' ? likedMovies : likedBooks;
  const type = activeTab === 'movies' ? 'movie' : 'book';

  const {
    data: batchMeta = {},
    isLoading: metaLoading,
    error: metaError,
  } = useQuery({
    queryKey: ['batchMeta', ids, type],
    queryFn: async () => {
      if (ids.length === 0) return {};
      const res = await api.get(`/api/items/batch?ids=${ids.join(',')}&type=${type}`);
      const map: Record<string, any> = {};
      (res.data.items || []).forEach((item: any) => {
        map[item.id] = {
          id: item.id,
          title: item.title,
          image: item.image,
          type: item.type,
        };
      });
      return map;
    },
    enabled: ids.length > 0,
    staleTime: 30 * 60 * 1000,
    gcTime: 60 * 60 * 1000,
    placeholderData: (prev) => prev,
  });

  const filteredIds = libraryQuery.trim()
    ? ids.filter(id => batchMeta[id]?.title?.toLowerCase().includes(libraryQuery.toLowerCase()))
    : ids;

  const addBook = useMutation({
    mutationFn: (bookId: string) => api.post(`/api/book/${bookId}/like`),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['likedBooks'] });
      setShowAddModal(false);
    },
  });
  
  const addMovie = useMutation({
    mutationFn: (movieId: string) => api.post(`/api/movie/${movieId}/like`),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['likedMovies'] });
      setShowAddModal(false);
    },
  });

  const isLoading = metaLoading;

  return (
    <div className="min-h-screen pb-20">
      <div className="px-4 pt-6 pb-2">
        <div className="flex items-center justify-between mb-4">
          <h1 className="text-2xl font-bold tracking-tight">Your Library</h1>
          <button
            onClick={() => setShowAddModal(true)}
            className="w-8 h-8 flex items-center justify-center rounded-full bg-primary text-primary-foreground"
          >
            <Plus className="w-4 h-4" />
          </button>
        </div>
        <div className="relative">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-muted-foreground" />
          <input
            type="text"
            placeholder={`Search your ${activeTab}...`}
            value={libraryQuery}
            onChange={(e) => setLibraryQuery(e.target.value)}
            className="w-full pl-9 pr-4 py-2 rounded-lg border border-border bg-background text-sm"
          />
        </div>
      </div>

      <div className="sticky top-0 z-10 bg-background/80 backdrop-blur-sm border-b border-border px-4 py-2">
        <div className="flex gap-2">
          <button
            onClick={() => handleTabChange('movies')}
            className={`px-4 py-2 rounded-full text-sm font-medium ${activeTab === 'movies' ? 'bg-primary text-primary-foreground' : 'bg-secondary'}`}
          >Movies</button>
          <button
            onClick={() => handleTabChange('books')}
            className={`px-4 py-2 rounded-full text-sm font-medium ${activeTab === 'books' ? 'bg-primary text-primary-foreground' : 'bg-secondary'}`}
          >Books</button>
        </div>
      </div>

      {isLoading ? (
        <div className="flex justify-center py-20"><Loader2 className="w-8 h-8 animate-spin" /></div>
      ) : metaError ? (
        <p className="text-destructive text-center py-20">Failed to load library</p>
      ) : filteredIds.length === 0 ? (
        <p className="text-muted-foreground text-center py-20">
          {libraryQuery.trim() ? 'No matches found.' : `No ${activeTab} in your library yet.`}
        </p>
      ) : (
        <div className="grid grid-cols-2 sm:grid-cols-[repeat(auto-fill,minmax(225px,1fr))] gap-2 p-4">
          {filteredIds.map(id => {
            const item = batchMeta[id];
            if (!item) return null;
            return item.type === 'movie' ? (
              <MovieCard key={id} movie={{ movie_id: id, title: item.title, poster_url: item.image }} />
            ) : (
              <BookCard key={id} book={{ book_id: id, title: item.title, image_url: item.image }} />
            );
          })}
        </div>
      )}

      {showAddModal && (
        <div className="fixed inset-0 z-50 bg-background/80 backdrop-blur-sm flex items-start justify-center pt-20">
          <div className="bg-background border border-border rounded-xl p-4 w-full max-w-md mx-4 shadow-lg">
            <div className="flex items-center justify-between mb-4">
              <h2 className="text-lg font-semibold">Add to Library</h2>
              <button onClick={() => setShowAddModal(false)}><X className="w-5 h-5" /></button>
            </div>
            <ItemSelector
              onSelect={(item) => {
                if (item.type === 'book') addBook.mutate(item.id);
                else addMovie.mutate(item.id);
              }}
            />
          </div>
        </div>
      )}
    </div>
  );
}