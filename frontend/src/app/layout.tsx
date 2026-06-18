'use client';

import { Inter } from 'next/font/google';
import './globals.css';
import { usePathname } from 'next/navigation';
import Providers from '@/components/providers';
import Sidebar from '@/components/sidebar';
import MobileNav from '@/components/mobile-nav';

const inter = Inter({ subsets: ['latin'] });

export default function RootLayout({ children }: { children: React.ReactNode }) {
  const pathname = usePathname();
  const hideNav = ['/login', '/register', '/onboarding', '/onboarding/picks'].includes(pathname);

  return (
    <html lang="en">
      <body className="pb-16 lg:pb-0">
        <div className="flex min-h-screen">
          {!hideNav && <Sidebar />}
          <main className={`flex-1 min-w-0 ${!hideNav ? 'pb-16 md:pb-0 md:ml-20' : ''}`}>
            <Providers>{children}</Providers>
          </main>
        </div>
        {!hideNav && <MobileNav />}
      </body>
    </html>
  );
}