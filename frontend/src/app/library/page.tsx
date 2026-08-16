'use client';

import { Suspense } from 'react';
import LibraryPage from '@/components/library-page';

export default function Page() {
  return (
    <Suspense fallback={<div className="min-h-screen flex justify-center items-center">Loading...</div>}>
      <LibraryPage />
    </Suspense>
  );
}