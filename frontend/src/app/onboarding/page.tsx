'use client';

import { useState, useEffect, useRef } from 'react';
import { useRouter } from 'next/navigation';
import { useQuery, useInfiniteQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { Loader2, Check, Plus, Search, Film, BookOpen } from 'lucide-react';
import api from '@/lib/api';
import MovieCard from '@/components/movie-card';
import BookCard from '@/components/book-card';
import { SelectableItem } from '@/components/item-selector';

type Tab = 'movies' | 'books';

export default function OnboardingPage() {
  const router = useRouter();
  const queryClient = useQueryClient();
  const [tab, setTab] = useState<Tab>('movies');
  const [selectedMovieGenre, setSelectedMovieGenre] = useState<string | null>(null);
  const [selectedBookGenre, setSelectedBookGenre] = useState<string | null>(null);
  const [selectedItems, setSelectedItems] = useState<SelectableItem[]>([]);
  const [searchQuery, setSearchQuery] = useState('');
  const [searchResults, setSearchResults] = useState<SelectableItem[]>([]);
  const [searchLoading, setSearchLoading] = useState(false);
  const debounceRef = useRef<ReturnType<typeof setTimeout> | null>(null);

  const isFirstRender = useRef(true);

  const { data: profile, isLoading: profileLoading } = useQuery({
    queryKey: ['userProfile'],
    queryFn: async () => (await api.get('/api/user/profile')).data,
  });
  useEffect(() => {
    if (profile?.onboarding_completed) router.replace('/');
  }, [profile, router]);

  const { data: movieGenres, isLoading: loadingMovieGenres } = useQuery({
    queryKey: ['genres', 'movie'],
    queryFn: async () => (await api.get('/api/genres?type=movie')).data.genres || [],
  });
  const { data: bookGenres, isLoading: loadingBookGenres } = useQuery({
    queryKey: ['genres', 'book'],
    queryFn: async () => (await api.get('/api/genres?type=book')).data.genres || [],
  });
  useEffect(() => {
    if (movieGenres?.length && !selectedMovieGenre) setSelectedMovieGenre(movieGenres[0]);
  }, [movieGenres]);
  useEffect(() => {
    if (bookGenres?.length && !selectedBookGenre) setSelectedBookGenre(bookGenres[0]);
  }, [bookGenres]);

  useEffect(() => {
    const saved = localStorage.getItem('onboardingPicks');
    if (saved) setSelectedItems(JSON.parse(saved));
  }, []);

  useEffect(() => {
    if (isFirstRender.current) {
      isFirstRender.current = false;
      return;
    }
    localStorage.setItem('onboardingPicks', JSON.stringify(selectedItems));
  }, [selectedItems]);

  useEffect(() => {
    const handlePopState = () => {
      const saved = localStorage.getItem('onboardingPicks');
      if (saved) setSelectedItems(JSON.parse(saved));
    };
    window.addEventListener('popstate', handlePopState);
    return () => window.removeEventListener('popstate', handlePopState);
  }, []);

  const currentGenres = tab === 'movies' ? movieGenres : bookGenres;
  const currentGenre = tab === 'movies' ? selectedMovieGenre : selectedBookGenre;
  const isLoadingGenres = tab === 'movies' ? loadingMovieGenres : loadingBookGenres;

  const {
    data,
    fetchNextPage,
    hasNextPage,
    isFetchingNextPage,
    isLoading: genreLoading,
  } = useInfiniteQuery({
    queryKey: ['onboardingGenre', tab, currentGenre],
    queryFn: async ({ pageParam = 0 }) => {
      if (!currentGenre) return [];
      const sort = tab === 'books' ? 'ratings_count:desc' : 'vote_count:desc';
      const type = tab === 'books' ? 'book' : 'movie';
      const limit = 30;
      const offset = pageParam * limit;
      const res = await api.get('/api/search', {
        params: { type, genre: currentGenre, sort, limit, offset },
      });
      return res.data.results || [];
    },
    initialPageParam: 0,
    getNextPageParam: (lastPage, allPages) => {
      if (lastPage.length < 30) return undefined;
      return allPages.length;
    },
    enabled: !!currentGenre && !searchQuery,
  });

  const genreItems = data?.pages.flatMap(page => page) ?? [];
  const selectedIds = new Set(selectedItems.map(i => i.id));

  useEffect(() => {
    if (debounceRef.current) clearTimeout(debounceRef.current);
    if (!searchQuery.trim()) {
      setSearchResults([]);
      return;
    }
    debounceRef.current = setTimeout(async () => {
      setSearchLoading(true);
      try {
        const res = await api.get(`/api/search?q=${encodeURIComponent(searchQuery)}&type=all`);
        const results = (res.data.results || [])
          .filter((item: any) => item.id)
          .map((item: any) => ({
            id: item.id,
            title: item.title,
            type: item.type,
            image: item.image,
            year: item.year,
            creator: item.type === 'book' ? item.authors : item.director,
          }));
        setSearchResults(results);
      } catch {
        setSearchResults([]);
      } finally {
        setSearchLoading(false);
      }
    }, 300);
    return () => { if (debounceRef.current) clearTimeout(debounceRef.current); };
  }, [searchQuery]);

  const toggleSelect = (item: SelectableItem) => {
    setSelectedItems(prev => {
      const exists = prev.some(i => i.id === item.id);
      if (exists) return prev.filter(i => i.id !== item.id);
      return [...prev, item];
    });
  };

  const completeOnboarding = useMutation({
    mutationFn: async () => {
      const bookIds = selectedItems.filter(i => i.type === 'book').map(i => i.id);
      const movieIds = selectedItems.filter(i => i.type === 'movie').map(i => i.id);
      for (const id of bookIds) await api.post(`/api/book/${id}/like`);
      for (const id of movieIds) await api.post(`/api/movie/${id}/like`);
      await api.post('/api/user/onboarding/complete', { completed: true });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['likedBooks'] });
      queryClient.invalidateQueries({ queryKey: ['likedMovies'] });
      queryClient.invalidateQueries({ queryKey: ['userProfile'] });
      router.push('/');
    },
  });

  const skipOnboarding = useMutation({
    mutationFn: async () => {
      await api.post('/api/user/onboarding/complete', { completed: true });
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['userProfile'] });
      router.push('/library');
    },
  });

  if (profileLoading || isLoadingGenres) {
    return <div className="min-h-screen flex justify-center items-center"><Loader2 className="w-8 h-8 animate-spin" /></div>;
  }

  return (
    <div className="min-h-screen pb-36 md:pb-20">
      <div className="px-4 pt-6 pb-2 flex justify-between items-start">
        <div>
          <h1 className="text-2xl font-bold">Choose your favorite movies or books</h1>
        </div>
      </div>

      <div className="px-4 pt-2 pb-2">
        <div className="relative">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-muted-foreground" />
          <input
            type="text"
            placeholder="Search books or movies..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            className="w-full pl-9 pr-4 py-2 rounded-lg border border-border bg-background text-sm"
          />
        </div>
        {searchLoading && (
          <div className="flex justify-center py-4">
            <Loader2 className="w-5 h-5 animate-spin text-muted-foreground" />
          </div>
        )}
        {!searchLoading && searchQuery.trim() !== '' && (
          <div className="space-y-2 mt-2">
            {searchResults.length === 0 ? (
              <p className="text-sm text-muted-foreground">Nothing found</p>
            ) : (
              searchResults.map((item) => (
                <div
                  key={`${item.type}-${item.id}`}
                  className="flex items-center justify-between border rounded-lg p-2 cursor-pointer hover:bg-secondary"
                  onClick={() => {
                    toggleSelect(item);
                    setSearchQuery('');
                  }}
                >
                  <div className="flex items-center gap-2 min-w-0">
                    {item.image ? (
                      <img src={item.image} alt={item.title} className="w-8 h-10 object-cover rounded flex-shrink-0" />
                    ) : (
                      <div className="w-8 h-10 bg-muted rounded flex-shrink-0 flex items-center justify-center">
                        {item.type === 'movie' ? <Film className="w-4 h-4" /> : <BookOpen className="w-4 h-4" />}
                      </div>
                    )}
                    <div>
                      <div className="font-medium text-sm truncate">{item.title}</div>
                      <div className="text-xs text-muted-foreground truncate">
                        {item.year && <span>{item.year}</span>}
                        {item.creator && <span>{item.year ? ' · ' : ''}{item.creator}</span>}
                      </div>
                    </div>
                  </div>
                  <Plus className="w-4 h-4 text-muted-foreground flex-shrink-0 ml-2" />
                </div>
              ))
            )}
          </div>
        )}
      </div>

      {!searchQuery && (
        <>
          <div className="flex gap-2 px-4 py-2">
            <button onClick={() => setTab('movies')} className={`px-4 py-2 rounded-lg text-sm font-medium ${tab === 'movies' ? 'bg-primary text-primary-foreground' : 'bg-secondary'}`}>Movies</button>
            <button onClick={() => setTab('books')} className={`px-4 py-2 rounded-lg text-sm font-medium ${tab === 'books' ? 'bg-primary text-primary-foreground' : 'bg-secondary'}`}>Books</button>
          </div>

          {currentGenres && currentGenres.length > 0 && (
            <div className="px-4 py-2 overflow-x-auto no-scrollbar">
              <div className="flex gap-2 min-w-max">
                {currentGenres.map((genre: string) => (
                  <button
                    key={genre}
                    onClick={() => tab === 'movies' ? setSelectedMovieGenre(genre) : setSelectedBookGenre(genre)}
                    className={`px-3 py-1 rounded-lg text-sm ${currentGenre === genre ? 'bg-primary text-primary-foreground' : 'bg-secondary'}`}
                  >
                    {genre.charAt(0).toUpperCase() + genre.slice(1)}
                  </button>
                ))}
              </div>
            </div>
          )}
  
          {selectedItems.length > 0 && (
            <div className="px-4 py-2">
              <button onClick={() => router.push('/onboarding/picks')} className="w-full py-2 bg-secondary rounded-lg text-sm flex justify-between items-center px-4">
                <span>Your current picks</span>
                <span className="bg-background/50 px-2 py-0.5 rounded-full text-xs">{selectedItems.length}</span>
              </button>
            </div>
          )}

          {genreLoading ? (
            <div className="flex justify-center py-20"><Loader2 className="w-8 h-8 animate-spin" /></div>
          ) : (
            <>
              <div className="grid grid-cols-2 sm:grid-cols-[repeat(auto-fill,minmax(225px,1fr))] gap-2 p-4">
                {genreItems.map((item: any) => {
                  const selectableItem: SelectableItem = {
                    id: item.id,
                    title: item.title,
                    type: item.type,
                    image: item.image,
                    year: item.year,
                    creator: item.type === 'book' ? item.authors : item.director,
                  };
                  const isSelected = selectedIds.has(item.id);
                  return (
                    <div
                      key={item.id}
                      className="relative cursor-pointer group"
                      onClick={() => toggleSelect(selectableItem)}
                    >
                      {item.type === 'movie' ? (
                        <MovieCard movie={{ movie_id: item.id, title: item.title, poster_url: item.image }}/>
                      ) : (
                        <BookCard book={{ book_id: item.id, title: item.title, image_url: item.image }}/>
                      )}
                      {isSelected && (
                        <div className="absolute inset-0 bg-black/60 flex items-center justify-center rounded-lg">
                          <Check className="w-22 h-22 text-white" />
                        </div>
                      )}
                      {!isSelected && (
                        <div className="absolute inset-0 bg-black/0 group-hover:bg-black/20 transition-colors rounded-lg" />
                      )}
                    </div>
                  );
                })}
              </div>

              {hasNextPage && (
                <div className="flex justify-center py-4">
                  <button onClick={() => fetchNextPage()} disabled={isFetchingNextPage} className="px-4 py-2 bg-secondary rounded-lg text-sm">
                    {isFetchingNextPage ? 'Loading...' : 'Load more'}
                  </button>
                </div>
              )}

              {genreItems.length === 0 && (
                <p className="text-center text-muted-foreground py-20">No {tab} found in this genre.</p>
              )}
            </>
          )}
        </>
      )}

      <div className="fixed bottom-0 left-0 right-0 p-4 flex justify-between items-center z-40 md:left-20 md:right-0 max-md:bottom-16">
        <button
          onClick={() => skipOnboarding.mutate()}
          disabled={skipOnboarding.isPending}
          className="px-12 py-3 rounded-lg border border-border bg-secondary text-foreground text-sm font-medium"
        >
          {skipOnboarding.isPending ? 'Skipping...' : 'Skip'}
        </button>
        <button
          onClick={() => completeOnboarding.mutate()}
          disabled={selectedItems.length === 0 || completeOnboarding.isPending}
          className="px-12 py-3 rounded-lg bg-primary text-primary-foreground text-sm font-medium disabled:opacity-50"
        >
          {completeOnboarding.isPending ? 'Saving...' : 'Complete'}
        </button>
      </div>
    </div>
  );
}