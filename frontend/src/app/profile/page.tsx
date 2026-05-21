'use client';

import { useQuery } from '@tanstack/react-query';
import { useRouter } from 'next/navigation';
import api from '@/lib/api';
import { HistoryEntry, SavedRecommendation } from '@/types';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';

export default function ProfilePage() {
  const router = useRouter();

  const { data: history } = useQuery<HistoryEntry[]>({
    queryKey: ['history'],
    queryFn: async () => (await api.get('/api/user/recommendations/history')).data.history,
  });

  const { data: saved } = useQuery<SavedRecommendation[]>({
    queryKey: ['saved'],
    queryFn: async () => (await api.get('/api/user/recommendations/saved')).data.saved,
  });

  const handleLogout = () => {
    localStorage.removeItem('token');
    document.cookie = 'token=; path=/; expires=Thu, 01 Jan 1970 00:00:00 UTC';
    router.push('/login');
  };

  return (
    <div className="container mx-auto p-4 pb-20 space-y-6">
      <h1 className="text-2xl font-bold">Profile</h1>

      <Card>
        <CardHeader><CardTitle>History</CardTitle></CardHeader>
        <CardContent>
          {history?.map((entry) => (
            <div key={entry.id} className="border-b py-2">
              <p className="text-sm">Task: {entry.task_id}</p>
              <p className="text-xs text-muted-foreground">Direction: {entry.direction}</p>
              {entry.result && <p className="text-xs">Result: {entry.result}</p>}
            </div>
          ))}
        </CardContent>
      </Card>

      <Card>
        <CardHeader><CardTitle>Saved</CardTitle></CardHeader>
        <CardContent>
          {saved?.map((item) => (
            <div key={item.id} className="border-b py-2">
              <p className="text-sm">{item.from_type}:{item.from_id} → {item.to_type}:{item.to_id}</p>
            </div>
          ))}
        </CardContent>
      </Card>

      <Button onClick={handleLogout} variant="destructive" className="w-full">
        Log out
      </Button>
    </div>
  );
}