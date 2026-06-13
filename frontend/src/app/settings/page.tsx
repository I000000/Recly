'use client';

import { useRouter } from 'next/navigation';
import { ArrowLeft, LogOut } from 'lucide-react';
import { Button } from '@/components/ui/button';

export default function SettingsPage() {
  const router = useRouter();

  const handleLogout = () => {
    localStorage.removeItem('token');
    document.cookie = 'token=; path=/; expires=Thu, 01 Jan 1970 00:00:00 UTC';
    router.push('/login');
  };

  return (
    <div className="min-h-screen pb-20">
      <div className="px-4 pt-6 pb-2 flex items-center gap-3">
        <h1 className="text-2xl font-bold tracking-tight">Settings</h1>
      </div>

      <div className="px-4 py-4 space-y-6">
        
        <div className="pt-4">
          <Button
            onClick={handleLogout}
            variant="destructive"
            className="w-full sm:w-auto"
          >
            <LogOut className="w-4 h-4 mr-2" />
            Log out
          </Button>
        </div>
      </div>
    </div>
  );
}