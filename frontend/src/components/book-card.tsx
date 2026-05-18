'use client';

import { useEffect, useState } from 'react';
import { loadBooks } from '@/lib/data';

export default function BookCard({ bookId }: { bookId: string }) {
  const [book, setBook] = useState<any>(null);
  useEffect(() => {
    loadBooks().then((books) => {
      const found = books.find((b: any) => b.book_id === bookId);
      if (found) setBook(found);
    });
  }, [bookId]);
  if (!book) return null;
  return (
    <div className="border rounded-lg overflow-hidden shadow-sm">
      {book.image_url && (
        <img src={book.image_url} alt={book.title} className="w-full h-48 object-cover" />
      )}
      <div className="p-3">
        <h3 className="font-semibold">{book.title}</h3>
        {book.average_rating && <p className="text-xs text-muted-foreground">★ {book.average_rating}</p>}
      </div>
    </div>
  );
}