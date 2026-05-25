'use client';

import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { Home, Library, User, Settings, Film } from 'lucide-react';

const links = [
  { href: '/', label: 'Home', icon: Home },
  { href: '/library', label: 'Library', icon: Library },
  { href: '/profile', label: 'Profile', icon: User },
];

export default function Sidebar() {
  const pathname = usePathname();

  return (
    <aside className="hidden md:flex flex-col fixed left-0 top-0 h-screen w-20 bg-background border-r border-border z-50">
      {/* Логотип */}
      <Link href="/" className="flex items-center justify-center h-20 border-b border-border">
        <Film className="w-7 h-7 text-primary" />
      </Link>

      {/* Навигация */}
      <nav className="flex-1 flex flex-col items-center gap-8 mt-8">
        {links.map(({ href, icon: Icon, label }) => (
          <Link
            key={href}
            href={href}
            className={`p-3 rounded-xl transition-colors ${
              pathname === href
                ? 'bg-primary/10 text-primary'
                : 'text-muted-foreground hover:bg-secondary hover:text-foreground'
            }`}
            title={label}
          >
            <Icon className="w-6 h-6" />
          </Link>
        ))}
      </nav>

      {/* Настройки внизу */}
      <Link
        href="/settings"
        className={`flex items-center justify-center h-20 border-t border-border ${
          pathname === '/settings'
            ? 'bg-primary/10 text-primary'
            : 'text-muted-foreground hover:bg-secondary hover:text-foreground'
        }`}
        title="Settings"
      >
        <Settings className="w-6 h-6" />
      </Link>
    </aside>
  );
}