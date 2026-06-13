'use client';

import { useRouter } from 'next/navigation';
import { useQueryClient } from '@tanstack/react-query';
import api from '@/lib/api';

export default function MovieCard({ movie, aspectRatio }: { movie: any; aspectRatio?: string }) {
  const router = useRouter();
  const queryClient = useQueryClient();
  const aspect = aspectRatio || '2/3';

  const handleClick = async () => {
    try {
      await api.post('/api/user/view', {
        item_id: movie.movie_id,
        item_type: 'movie',
      });
      queryClient.invalidateQueries({ queryKey: ['views'] });
    } catch (e) {
      console.error('Failed to record view', e);
    }
    router.push(`/movie/${movie.movie_id}`);
  };

  return (
    <div onClick={handleClick} className="block cursor-pointer">
      <div
        className="relative border rounded-xl overflow-hidden shadow-md group"
        style={{ aspectRatio: aspect }}
      >
        {movie.poster_url && (
          <img
            src={movie.poster_url}
            alt={movie.title}
            className="absolute inset-0 w-full h-full object-cover transition-transform duration-300 group-hover:scale-105"
            onError={(e) => {
              (e.target as HTMLImageElement).src = '/placeholder.png';
            }}
          />
        )}
        <div className="absolute bottom-0 left-0 right-0 bg-gradient-to-t from-black/70 via-black/20 to-transparent p-3 pt-22">
          <h3 className="text-white text-sm font-semibold leading-snug line-clamp-2 drop-shadow-md">
            {movie.title}
          </h3>
        </div>
      </div>
    </div>
  );
}