import Link from 'next/link';

export default function MovieCard({ movie }: { movie: any }) {
  return (
    <Link href={`/movie/${movie.movie_id}`} className="block">
      <div className="relative border rounded-xl overflow-hidden shadow-md aspect-[2/3] group">
        {movie.poster_url && (
          <img
            src={movie.poster_url}
            alt={movie.title}
            className="absolute inset-0 w-full h-full object-cover transition-transform duration-300 group-hover:scale-105"
            referrerPolicy="no-referrer"
            crossOrigin="anonymous"
            onError={(e) => {
              (e.target as HTMLImageElement).src = '/placeholder.png';
            }}
          />
        )}
        <div className="absolute bottom-0 left-0 right-0 bg-gradient-to-t from-black/80 via-black/40 to-transparent p-3 pt-12">
          <h3 className="text-white text-sm font-semibold leading-snug line-clamp-2 drop-shadow-md">
            {movie.title}
          </h3>
        </div>
      </div>
    </Link>
  );
}