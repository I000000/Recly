'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { useQuery, useInfiniteQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { Loader2, Check, Plus } from 'lucide-react';
import api from '@/lib/api';
import MovieCard from '@/components/movie-card';
import BookCard from '@/components/book-card';
import ItemSelector, { SelectableItem } from '@/components/item-selector';

type Tab = 'movies' | 'books';

export default function OnboardingPage() {
  const router = useRouter();
  const queryClient = useQueryClient();
  const [tab, setTab] = useState<Tab>('movies');
  const [selectedMovieGenre, setSelectedMovieGenre] = useState<string | null>(null);
  const [selectedBookGenre, setSelectedBookGenre] = useState<string | null>(null);
  const [selectedItems, setSelectedItems] = useState<SelectableItem[]>([]);
  const [searchQuery, setSearchQuery] = useState('');

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
    localStorage.setItem('onboardingPicks', JSON.stringify(selectedItems));
  }, [selectedItems]);

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

  const toggleSelect = (item: SelectableItem) => {
    setSelectedItems(prev => {
      const exists = prev.some(i => i.id === item.id);
      if (exists) return prev.filter(i => i.id !== item.id);
      return [...prev, item];
    });
  };

  const completeOnboarding = useMutation({ /* как раньше */ });
  const skipOnboarding = useMutation({ /* как раньше */ });

  if (profileLoading || isLoadingGenres) {
    return <div className="min-h-screen flex justify-center items-center"><Loader2 className="w-8 h-8 animate-spin" /></div>;
  }

  return (
    <div className="min-h-screen pb-24">
      <div className="px-4 pt-6 pb-2 flex justify-between items-start">
        <div>
          <h1 className="text-2xl font-bold">Choose your favorite movies or books</h1>
        </div>
      </div>

      <div className="p-4 px-4">
        <ItemSelector onSelect={toggleSelect} searchQuery={searchQuery} setSearchQuery={setSearchQuery} expandResults={!!searchQuery}/>
      </div>

      {!searchQuery && (
        <>
          <div className="flex gap-2 px-4">
            <button onClick={() => setTab('movies')} className={`px-4 py-2 rounded-full text-sm ${tab === 'movies' ? 'bg-primary text-primary-foreground' : 'bg-secondary'}`}>Movies</button>
            <button onClick={() => setTab('books')} className={`px-4 py-2 rounded-full text-sm ${tab === 'books' ? 'bg-primary text-primary-foreground' : 'bg-secondary'}`}>Books</button>
          </div>

          {currentGenres && currentGenres.length > 0 && (
            <div className="px-4 py-4 overflow-x-auto no-scrollbar">
              <div className="flex gap-2 min-w-max">
                {currentGenres.map((genre: string) => (
                  <button
                    key={genre}
                    onClick={() => tab === 'movies' ? setSelectedMovieGenre(genre) : setSelectedBookGenre(genre)}
                    className={`px-3 py-1 rounded-full text-sm ${currentGenre === genre ? 'bg-primary text-primary-foreground' : 'bg-secondary'}`}
                  >
                    {genre.charAt(0).toUpperCase() + genre.slice(1)}
                  </button>
                ))}
              </div>
            </div>
          )}

          {selectedItems.length > 0 && (
            <div className="px-4">
              <button onClick={() => router.push('/onboarding/picks')} className="w-full py-2 bg-secondary rounded-lg text-sm">Your current picks ({selectedItems.length})</button>
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

      {selectedItems.length > 0 && !searchQuery && (
        <div className="fixed bottom-0 left-0 right-0 p-4 bg-background border-t">
          <button onClick={() => completeOnboarding.mutate()} disabled={completeOnboarding.isPending} className="w-full py-2 bg-primary text-primary-foreground rounded-lg font-medium">
            {completeOnboarding.isPending ? 'Saving...' : 'Complete onboarding'}
          </button>
        </div>
      )}
    </div>
  );
}