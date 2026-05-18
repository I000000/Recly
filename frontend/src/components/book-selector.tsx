export default function BookCard({ book }: { book: any }) {
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