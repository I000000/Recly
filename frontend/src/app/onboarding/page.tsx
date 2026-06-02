'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import api from '@/lib/api';
import ItemSelector from '@/components/item-selector';
import type { SelectableItem } from '@/components/item-selector';
import { Button } from '@/components/ui/button';
import { X } from 'lucide-react';

type PickedItem = {
  id: string;
  title: string;
  type: 'book' | 'movie';
  image?: string;
};

export default function OnboardingPage() {
  const [selected, setSelected] = useState<PickedItem[]>([]);
  const router = useRouter();

  const handleSelect = (item: SelectableItem) => {
    if (selected.length >= 5) return;
    if (selected.find((s) => s.id === item.id && s.type === item.type)) return;
    setSelected((prev) => [
      ...prev,
      {
        id: item.id,
        title: item.title,
        type: item.type,
        image: item.image,
      },
    ]);
  };

  const handleRemove = (id: string, type: 'book' | 'movie') => {
    setSelected((prev) => prev.filter((item) => !(item.id === id && item.type === type)));
  };

  const finish = async () => {
    for (const item of selected) {
      if (item.type === 'book') {
        await api.post(`/api/book/${item.id}/like`);
      } else {
        await api.post(`/api/movie/${item.id}/like`);
      }
    }
    router.push('/');
  };

  return (
    <div className="container mx-auto p-4 pb-20">
      <h1 className="text-2xl font-bold mb-4">Pick your favorites</h1>
      <p className="text-muted-foreground mb-4">Choose at least 3 books or movies you love</p>

      <ItemSelector onSelect={handleSelect} />

      {selected.length > 0 && (
        <div className="mt-4 space-y-2">
          {selected.map((item) => (
            <div
              key={`${item.type}-${item.id}`}
              className="flex items-center justify-between bg-secondary rounded-lg px-3 py-2"
            >
              <div className="flex items-center gap-2">
                {item.image && (
                  <img
                    src={item.image}
                    alt={item.title}
                    className="w-8 h-10 object-cover rounded"
                  />
                )}
                <div>
                  <span className="text-sm font-medium">{item.title}</span>
                  <span className="ml-2 text-xs text-muted-foreground">
                    {item.type === 'book' ? '📖' : '🎬'}
                  </span>
                </div>
              </div>
              <button
                onClick={() => handleRemove(item.id, item.type)}
                className="text-muted-foreground hover:text-destructive"
              >
                <X className="w-4 h-4" />
              </button>
            </div>
          ))}
        </div>
      )}

      <div className="mt-4">
        <p className="text-sm text-muted-foreground">Selected: {selected.length}</p>
        <Button onClick={finish} disabled={selected.length < 3} className="mt-2 w-full">
          Continue
        </Button>
      </div>
    </div>
  );
}