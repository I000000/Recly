import type { Metadata } from 'next';
import { Inter } from 'next/font/google';
import './globals.css';
import Providers from '@/components/providers';

const inter = Inter({ subsets: ['latin'] });

export const metadata: Metadata = {
  title: 'Recly',
  description: 'Discover movies and books recommendations just for you',
};

import MobileNav from '@/components/mobile-nav';

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en">
      <body className="pb-16 lg:pb-0"> {/* отступ для мобильной панели */}
        <Providers>{children}</Providers>
        <MobileNav />
      </body>
    </html>
  );
}