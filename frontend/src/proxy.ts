import { NextResponse } from 'next/server';
import type { NextRequest } from 'next/server';

export function proxy(request: NextRequest) {
  const token = request.cookies.get('token')?.value;

  // Разрешаем доступ к статическим данным без авторизации
  if (request.nextUrl.pathname.startsWith('/data/')) {
    return NextResponse.next();
  }

  // Публичные страницы, доступные без токена
  const isPublicPage =
    request.nextUrl.pathname.startsWith('/login') ||
    request.nextUrl.pathname.startsWith('/register') ||
    request.nextUrl.pathname.startsWith('/onboarding');

  // Если нет токена и страница не публичная → отправляем на /login
  if (!token && !isPublicPage) {
    return NextResponse.redirect(new URL('/login', request.url));
  }

  // Если токен есть и пользователь на странице входа или регистрации → на главную
  if (token && (request.nextUrl.pathname.startsWith('/login') || request.nextUrl.pathname.startsWith('/register'))) {
    return NextResponse.redirect(new URL('/', request.url));
  }

  return NextResponse.next();
}

export const config = {
  matcher: ['/((?!api|_next|static|public|favicon.ico).*)'],
};