'use client';

import { Direction } from '@/types';

const directions: { value: Direction; label: string }[] = [
  { value: 'book_to_movie', label: 'Books → Movies' },
  { value: 'movie_to_book', label: 'Movies → Books' },
  { value: 'book_to_book', label: 'Books → Books' },
  { value: 'movie_to_movie', label: 'Movies → Movies' },
];

export default function DirectionSwitcher({
  value,
  onChange,
}: {
  value: Direction;
  onChange: (dir: Direction) => void;
}) {
  return (
    <div className="flex gap-2">
      {directions.map((dir) => (
        <button
          key={dir.value}
          onClick={() => onChange(dir.value)}
          className={`px-3 py-1 rounded ${
            value === dir.value ? 'bg-primary text-white' : 'bg-secondary'
          }`}
        >
          {dir.label}
        </button>
      ))}
    </div>
  );
}