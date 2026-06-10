'use client';

import { useParams, useRouter } from 'next/navigation';
import { useQuery } from '@tanstack/react-query';
import { ArrowLeft, Loader2 } from 'lucide-react';
import api from '@/lib/api';
import MovieCard from '@/components/movie-card';

export default function MovieGenrePage() {
  const { genreName } = useParams() as { genreName: string };
  const router = useRouter();
  const decodedGenre = decodeURIComponent(genreName);

  const { data, isLoading, isError } = useQuery({
    queryKey: ['movieGenre', decodedGenre],
    queryFn: async () => {
      const res = await api.get('/api/search', {
        params: {
          type: 'movie',
          genre: decodedGenre,
          sort: 'vote_count:desc',
          limit: 30,
        },
      });
      return res.data.results || [];
    },
  });

  if (isLoading) {
    return (
      <div className="min-h-screen flex justify-center items-center">
        <Loader2 className="w-8 h-8 animate-spin text-muted-foreground" />
      </div>
    );
  }

  if (isError) {
    return <div className="p-8 text-destructive">Failed to load movies for this genre.</div>;
  }

  return (
    <div className="min-h-screen pb-20">
      <div className="px-4 pt-6 pb-2 flex items-center gap-3">
        <button onClick={() => router.back()} className="p-2 rounded-full hover:bg-secondary">
          <ArrowLeft className="w-5 h-5" />
        </button>
        <h1 className="text-2xl font-bold tracking-tight">
          Popular Movies in {decodedGenre.toUpperCase()}
        </h1>
      </div>

      {data.length === 0 ? (
        <p className="text-muted-foreground text-center py-20">
          No movies found in {decodedGenre.toUpperCase()}.
        </p>
      ) : (
        <div className="grid grid-cols-2 sm:grid-cols-[repeat(auto-fill,minmax(225px,1fr))] gap-2 p-4">
          {data.map((item: any) => (
            <MovieCard
              key={item.id}
              movie={{ movie_id: item.id, title: item.title, poster_url: item.image }}
              aspectRatio="2/3"
            />
          ))}
        </div>
      )}
    </div>
  );
}