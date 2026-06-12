'use client';

import { useRouter } from 'next/navigation';
import { Trash2, ArrowLeft } from 'lucide-react';
import { useEffect, useState } from 'react';
import { SelectableItem } from '@/components/item-selector';
import MovieCard from '@/components/movie-card';
import BookCard from '@/components/book-card';

export default function PicksPage() {
  const router = useRouter();
  const [items, setItems] = useState<SelectableItem[]>([]);

  useEffect(() => {
    const saved = localStorage.getItem('onboardingPicks');
    if (saved) setItems(JSON.parse(saved));
  }, []);

  const removeItem = (id: string) => {
    const updated = items.filter(i => i.id !== id);
    setItems(updated);
    localStorage.setItem('onboardingPicks', JSON.stringify(updated));
  };

  return (
    <div className="min-h-screen pb-20">
      <div className="px-4 pt-6 pb-2 flex items-center gap-3">
        <button onClick={() => router.back()} className="p-2 rounded-full hover:bg-secondary">
          <ArrowLeft className="w-5 h-5" />
        </button>
        <h1 className="text-2xl font-bold tracking-tight">Your current picks</h1>
      </div>
      {items.length === 0 ? (
        <p className="text-center py-20 text-muted-foreground">No items selected</p>
      ) : (
        <div className="grid grid-cols-2 sm:grid-cols-[repeat(auto-fill,minmax(225px,1fr))] gap-2 p-4">
          {items.map(item => (
            <div key={item.id} className="relative">
              {item.type === 'movie' ? (
                <MovieCard movie={{ movie_id: item.id, title: item.title, poster_url: item.image }} />
              ) : (
                <BookCard book={{ book_id: item.id, title: item.title, image_url: item.image }} aspectRatio="2/3" />
              )}
              <button
                onClick={() => removeItem(item.id)}
                className="absolute top-2 right-2 p-2 rounded-full bg-destructive/80 text-white hover:bg-destructive"
              >
                <Trash2 className="w-5 h-5" />
              </button>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}