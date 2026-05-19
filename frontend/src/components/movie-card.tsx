export default function MovieCard({ movie }: { movie: any }) {
  return (
    <div className="relative border rounded-lg overflow-hidden shadow-sm aspect-[2/3]">
      {movie.poster_url && (
        <img
          src={movie.poster_url}
          alt={movie.title}
          className="absolute inset-0 w-full h-full object-cover"
          referrerPolicy="no-referrer"
          crossOrigin="anonymous"
          onError={(e) => {
            (e.target as HTMLImageElement).src = '/placeholder.png';
          }}
        />
      )}
      <div className="absolute bottom-0 left-0 right-0 bg-gradient-to-t from-black/70 via-black/20 to-transparent p-3 pt-22">
        <h3 className="text-white text-sm font-bold leading-tight line-clamp-2">
          {movie.title}
        </h3>
      </div>
    </div>
  );
}