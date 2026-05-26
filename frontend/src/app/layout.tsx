import type { Metadata } from 'next';
import { Inter } from 'next/font/google';
import './globals.css';
import Providers from '@/components/providers';
import Sidebar from '@/components/sidebar';
import MobileNav from '@/components/mobile-nav';

const inter = Inter({ subsets: ['latin'] });

export const metadata: Metadata = {
  title: 'Recly',
  description: 'Discover movies and books recommendations just for you',
};

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en" className="overflow-x-hidden">
      <body className="overflow-x-hidden pb-16 lg:pb-0">
        <div className="flex min-h-screen">
          <Sidebar />
          <main className="flex-1 pb-16 md:pb-0 md:ml-20 max-w-[100vw] overflow-x-hidden">
            <Providers>{children}</Providers>
          </main>
        </div>
        <MobileNav />
      </body>
    </html>
  );
}