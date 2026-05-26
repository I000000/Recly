import Link from 'next/link';

export default function BookCard({ book, aspectRatio }: { book: any; aspectRatio?: string }) {
  const aspect = aspectRatio || '3/4';
  return (
    <Link href={`/book/${book.book_id}`} className="block">
      <div
        className="relative border rounded-xl overflow-hidden shadow-md group"
        style={{ aspectRatio: aspect }}
      >
        {book.image_url && (
          <img
            src={book.image_url}
            alt={book.title}
            className="absolute inset-0 w-full h-full object-cover transition-transform duration-300 group-hover:scale-105"
            onError={(e) => {
              (e.target as HTMLImageElement).src = '/placeholder.png';
            }}
          />
        )}
        <div className="absolute bottom-0 left-0 right-0 bg-gradient-to-t from-black/70 via-black/20 to-transparent p-3 pt-22">
          <h3 className="text-white text-sm font-semibold leading-snug line-clamp-2 drop-shadow-md">
            {book.title}
          </h3>
        </div>
      </div>
    </Link>
  );
}