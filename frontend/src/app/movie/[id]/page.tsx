'use client';
import { useEffect, useState } from 'react';
import { useParams } from 'next/navigation';
import { useQuery } from '@tanstack/react-query';
import api from '@/lib/api';
import MovieCard from '@/components/movie-card';
import BookCard from '@/components/book-card';

export default function MoviePage() {
  const params = useParams();
  const movieId = params.id as string;
  const [movie, setMovie] = useState<any>(null);

  const getYear = (dateStr: string | undefined): string => {
    if (!dateStr) return '';
    const match = dateStr.match(/^\d{4}/);
    return match ? match[0] : '';
  };

  // Загружаем информацию о фильме
  useEffect(() => {
    api.get(`/api/items/batch?ids=${movieId}&type=movie`)
      .then(res => {
        const items = res.data.items || [];
        if (items.length > 0) setMovie(items[0]);
      })
      .catch(() => setMovie(null));
  }, [movieId]);

  // Похожие фильмы
  const { data: similarMovies } = useQuery<any[]>({
    queryKey: ['similarMovies', movieId],
    queryFn: async () => {
      const res = await api.post('/api/recommend', {
        selected_ids: [`movie_${movieId}`],
        direction: 'movie_to_movie',
        weights: { genre: 0.3, text: 0.4, image: 0.3 },
      });
      let result = await api.get(`/api/result/${res.data.task_id}`);
      while (result.data.status === 'pending') {
        await new Promise((r) => setTimeout(r, 1000));
        result = await api.get(`/api/result/${res.data.task_id}`);
      }
      const ids = result.data.movies || [];
      if (ids.length === 0) return [];
      const batch = await api.get(`/api/items/batch?ids=${ids.join(',')}&type=movie`);
      return batch.data.items || [];
    },
    enabled: !!movieId,
  });

  // Похожие книги
  const { data: similarBooks } = useQuery<any[]>({
    queryKey: ['similarBooks', movieId],
    queryFn: async () => {
      const res = await api.post('/api/recommend', {
        selected_ids: [`movie_${movieId}`],
        direction: 'movie_to_book',
        weights: { genre: 0.3, text: 0.4, image: 0.3 },
      });
      let result = await api.get(`/api/result/${res.data.task_id}`);
      while (result.data.status === 'pending') {
        await new Promise((r) => setTimeout(r, 1000));
        result = await api.get(`/api/result/${res.data.task_id}`);
      }
      const ids = result.data.movies || [];
      if (ids.length === 0) return [];
      const batch = await api.get(`/api/items/batch?ids=${ids.join(',')}&type=book`);
      return batch.data.items || [];
    },
    enabled: !!movieId,
  });

  if (!movie) return <div className="p-8 text-center">Loading...</div>;

  return (
    <div className="min-h-screen pb-20 px-4">
      <div className="pt-6">
        {movie.image_url && (
          <img src={movie.image_url} alt={movie.title} className="w-full h-64 object-cover rounded-xl mb-4" />
        )}
        <h1 className="text-2xl font-bold">{movie.title}</h1>
        {movie.director && <p className="text-muted-foreground">Director: {movie.director}</p>}
        {movie.cast && <p className="text-sm text-muted-foreground">Cast: {movie.cast}</p>}
        {movie.genres && (
          <div className="flex flex-wrap gap-1 mt-2">
            {movie.genres.map((g: string) => (
              <span key={g} className="px-2 py-0.5 bg-secondary rounded-full text-xs">{g}</span>
            ))}
          </div>
        )}
        <div className="flex items-center gap-4 mt-2">
          {movie.rating > 0 && <p>★ {movie.rating}</p>}
          {movie.year > 0 && <p className="text-muted-foreground">{movie.year}</p>}
          {movie.runtime > 0 && <p className="text-muted-foreground">{movie.runtime} min</p>}
        </div>
        {movie.description && <p className="mt-4 text-sm">{movie.description}</p>}
      </div>

      {similarMovies && similarMovies.length > 0 && (
        <div className="mt-8">
          <h2 className="text-xl font-semibold mb-4">Similar Movies</h2>
          <div className="flex gap-2 overflow-x-auto pb-2 no-scrollbar">
            {similarMovies.slice(0, 10).map((item: any) => (
              <div key={item.id} className="flex-shrink-0 w-28">
                <MovieCard movie={{ movie_id: item.id, title: item.title, poster_url: item.image }} />
              </div>
            ))}
          </div>
        </div>
      )}

      {similarBooks && similarBooks.length > 0 && (
        <div className="mt-8">
          <h2 className="text-xl font-semibold mb-4">Books You Might Like</h2>
          <div className="flex gap-2 overflow-x-auto pb-2 no-scrollbar">
            {similarBooks.slice(0, 10).map((item: any) => (
              <div key={item.id} className="flex-shrink-0 w-28">
                <BookCard book={{ book_id: item.id, title: item.title, image_url: item.image }} />
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  );
}