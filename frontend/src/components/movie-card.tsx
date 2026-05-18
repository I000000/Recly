export default function MovieCard({ movie }: { movie: any }) {
  return (
    <div className="border rounded-lg overflow-hidden shadow-sm">
      {movie.poster_url && (
        <img
          src={`/posters/${movie.movie_id}.jpg`}
          alt={movie.title}
          className="w-full aspect-[2/3] object-cover rounded-t-lg"
          onError={(e) => {
            (e.target as HTMLImageElement).src = movie.poster_url || '/placeholder.png';
          }}
        />
      )}
      <div className="p-3">
        <h3 className="font-semibold">{movie.title}</h3>
        {movie.vote_average && <p className="text-xs text-muted-foreground">★ {movie.vote_average}</p>}
      </div>
    </div>
  );
}