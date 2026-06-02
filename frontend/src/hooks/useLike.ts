import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import api from '@/lib/api';

type ItemType = 'book' | 'movie';

export function useLike(itemType: ItemType, itemId: string) {
  const queryClient = useQueryClient();
  const queryKey = itemType === 'book' ? ['likedBooks'] : ['likedMovies'];
  const entity = itemType === 'book' ? 'book' : 'movie';

  const { data: likedRaw } = useQuery<string[]>({
    queryKey,
    queryFn: async () => {
      const res = await api.get(`/api/user/library/${itemType}s`);
      return res.data[`${itemType}s`].map((b: any) => b[`${itemType}_id`]);
    },
    staleTime: 1000 * 60 * 30,
  });

  const liked = likedRaw ?? [];
  const isLiked = liked.includes(itemId);

  const like = useMutation({
    mutationFn: () => api.post(`/api/${entity}/${itemId}/like`),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey });
      queryClient.invalidateQueries({ queryKey: ['homeRecommendations', itemType === 'book' ? 'books' : 'movies'] });
    },
  });

  const unlike = useMutation({
    mutationFn: () => api.delete(`/api/${entity}/${itemId}/like`),
    onSuccess: () => {
      queryClient.setQueryData<string[]>(queryKey, (old) =>
        (old || []).filter(id => id !== itemId)
      );
      queryClient.invalidateQueries({ queryKey });
      queryClient.invalidateQueries({ queryKey: ['homeRecommendations', itemType === 'book' ? 'books' : 'movies'] });
    },
  });

  const toggleLike = () => {
    if (isLiked) {
      unlike.mutate();
    } else {
      like.mutate();
    }
  };

  return {
    isLiked,
    toggleLike,
    isPending: like.isPending || unlike.isPending,
  };
}