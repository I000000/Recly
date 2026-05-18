import { Book, Movie } from '@/types';

let booksCache: Book[] | null = null;
let moviesCache: Movie[] | null = null;

export async function loadBooks(): Promise<Book[]> {
  if (booksCache) return booksCache;
  const res = await fetch('/data/books.json');
  booksCache = await res.json();
  return booksCache!;
}

export async function loadMovies(): Promise<Movie[]> {
  if (moviesCache) return moviesCache;
  const res = await fetch('/data/movies.json');
  moviesCache = await res.json();
  return moviesCache!;
}