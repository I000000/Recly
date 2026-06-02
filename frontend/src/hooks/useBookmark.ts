import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import api from '@/lib/api';

export function useBookmark(itemType: 'book' | 'movie', itemId: string) {
  const queryClient = useQueryClient();

  const { data: savedRaw = [] } = useQuery<any[]>({
    queryKey: ['saved'],
    queryFn: async () => (await api.get('/api/user/saved-items')).data.saved,
  });
  const saved = savedRaw ?? [];

  const savedEntry = saved.find(
    (s: any) => s.item_type === itemType && s.item_id === itemId
  );
  const isBookmarked = !!savedEntry;

  const save = useMutation({
    mutationFn: () =>
      api.post('/api/user/saved-items', {
        item_type: itemType,
        item_id: itemId,
      }),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['saved'] }),
  });

  const unsave = useMutation({
    mutationFn: () => api.delete(`/api/user/saved-items/${savedEntry?.id}`),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['saved'] }),
  });

  const toggleBookmark = () => {
    if (isBookmarked) {
      unsave.mutate();
    } else {
      save.mutate();
    }
  };

  return { isBookmarked, toggleBookmark, isPending: save.isPending || unsave.isPending };
}