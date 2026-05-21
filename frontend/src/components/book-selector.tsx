'use client';

import { useState, useEffect, useRef } from 'react';
import { Plus } from 'lucide-react';
import api from '@/lib/api';
import { Input } from '@/components/ui/input';

export type SelectableItem = {
  id: string;
  title: string;
  type: 'book' | 'movie';
  image?: string;
};

export default function BookSelector({ onSelect }: { onSelect: (item: SelectableItem) => void }) {
  const [query, setQuery] = useState('');
  const [items, setItems] = useState<SelectableItem[]>([]);
  const [loading, setLoading] = useState(false);
  const debounceRef = useRef<ReturnType<typeof setTimeout> | null>(null);

  useEffect(() => {
    if (debounceRef.current) clearTimeout(debounceRef.current);
    if (!query.trim()) {
      setItems([]);
      return;
    }
    debounceRef.current = setTimeout(async () => {
      setLoading(true);
      try {
        const res = await api.get(`/api/search?q=${encodeURIComponent(query)}&type=all`);
        const results: SelectableItem[] = (res.data.results || [])
          .filter((item: any) => item.id)
          .map((item: any) => ({
            id: item.id,
            title: item.title,
            type: item.type,
            image: item.image,
          }));
        setItems(results);
      } catch {
        setItems([]);
      } finally {
        setLoading(false);
      }
    }, 300);
    return () => { if (debounceRef.current) clearTimeout(debounceRef.current); };
  }, [query]);

  return (
    <div>
      <Input
        value={query}
        onChange={(e) => setQuery(e.target.value)}
        placeholder="Search books or movies..."
        className="mb-2"
      />
      {loading && (
        <div className="flex justify-center py-4">
          <div className="w-5 h-5 border-2 border-primary border-t-transparent rounded-full animate-spin" />
        </div>
      )}
      <div className="space-y-2 max-h-60 overflow-y-auto">
        {items.map((item) => (
          <div
            key={`${item.type}-${item.id}`}
            className="flex items-center justify-between border rounded p-2 cursor-pointer hover:bg-secondary active:bg-secondary touch-manipulation"
            onClick={() => onSelect(item)}
          >
            <div className="flex items-center gap-2">
              {item.image && (
                <img src={item.image} alt={item.title} className="w-8 h-10 object-cover rounded" />
              )}
              <div>
                <span className="font-medium text-sm">{item.title}</span>
                <span className="ml-2 text-xs text-muted-foreground">
                  {item.type === 'book' ? '📖' : '🎬'}
                </span>
              </div>
            </div>
            <Plus className="w-4 h-4 text-muted-foreground flex-shrink-0" />
          </div>
        ))}
        {query && !loading && items.length === 0 && (
          <p className="text-sm text-muted-foreground">Nothing found</p>
        )}
      </div>
    </div>
  );
}