'use client';

import { useState, useEffect, useRef } from 'react';
import { Film, BookOpen, Plus } from 'lucide-react';
import api from '@/lib/api';
import { Input } from '@/components/ui/input';

export type SelectableItem = {
  id: string;
  title: string;
  type: 'book' | 'movie';
  image?: string;
  year?: number;
  creator?: string;
};

type ItemSelectorProps = {
  onSelect: (item: SelectableItem) => void;
  searchQuery?: string;
  setSearchQuery?: (query: string) => void;
  expandResults?: boolean;
};

export default function ItemSelector({ onSelect, searchQuery: externalQuery, setSearchQuery: externalSetQuery, expandResults }: ItemSelectorProps) {
  const [internalQuery, setInternalQuery] = useState('');
  const isControlled = externalQuery !== undefined && externalSetQuery !== undefined;
  const query = isControlled ? externalQuery : internalQuery;
  const setQuery = isControlled ? externalSetQuery : setInternalQuery;

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
            year: item.year,
            creator: item.type === 'book' ? item.authors : item.director,
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
    <div className={expandResults ? "flex flex-col h-full" : ""}>
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
      <div className={expandResults ? "flex-1 overflow-y-auto" : "space-y-2 max-h-60 overflow-y-auto"}>
        {items.map((item) => (
          <div
            key={`${item.type}-${item.id}`}
            className="flex items-center justify-between border rounded p-2 cursor-pointer hover:bg-secondary active:bg-secondary touch-manipulation"
            onClick={() => onSelect(item)}
          >
            <div className="flex items-center gap-2 min-w-0">
              {item.image ? (
                <img src={item.image} alt={item.title} className="w-8 h-10 object-cover rounded flex-shrink-0" />
              ) : (
                <div className="w-8 h-10 bg-muted rounded flex-shrink-0 flex items-center justify-center">
                  {item.type === 'movie' ? <Film className="w-4 h-4 text-muted-foreground" /> : <BookOpen className="w-4 h-4 text-muted-foreground" />}
                </div>
              )}
              <div className="min-w-0">
                <div className="flex items-center gap-1">
                  {item.type === 'movie' ? (
                    <Film className="w-3.5 h-3.5 text-muted-foreground flex-shrink-0" />
                  ) : (
                    <BookOpen className="w-3.5 h-3.5 text-muted-foreground flex-shrink-0" />
                  )}
                  <span className="font-medium text-sm truncate">{item.title}</span>
                </div>
                <div className="text-xs text-muted-foreground truncate">
                  {item.year && <span>{item.year}</span>}
                  {item.creator && <span>{item.year ? ' · ' : ''}{item.creator}</span>}
                </div>
              </div>
            </div>
            <Plus className="w-4 h-4 text-muted-foreground flex-shrink-0 ml-2" />
          </div>
        ))}
        {query && !loading && items.length === 0 && (
          <p className="text-sm text-muted-foreground">Nothing found</p>
        )}
      </div>
    </div>
  );
}