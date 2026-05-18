'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import api from '@/lib/api';
import BookSelector from '@/components/book-selector';

export default function OnboardingPage() {
  const [selected, setSelected] = useState<{ id: string; type: 'book' | 'movie' }[]>([]);
  const router = useRouter();

  const handleSelect = (item: any) => {
    if (selected.length >= 5) return;
    setSelected([...selected, { id: item.type === 'book' ? item.book_id : item.movie_id, type: item.type }]);
  };

  const finish = async () => {
    // Добавляем каждую выбранную книгу/фильм в библиотеку
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
    <div className="container mx-auto p-4">
      <h1 className="text-2xl font-bold mb-4">Выберите любимые книги и фильмы (минимум 3)</h1>
      <BookSelector onSelect={handleSelect} />
      <div className="mt-4">
        <h2 className="font-semibold">Выбрано: {selected.length}</h2>
        <Button onClick={finish} disabled={selected.length < 3} className="mt-2">
          Продолжить
        </Button>
      </div>
    </div>
  );
}