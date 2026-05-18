export interface Book {
    book_id: string;
    title: string;
    image_url?: string;
    average_rating?: number;
  }
  
  export interface Movie {
    movie_id: string;
    title: string;
    poster_url?: string;
    vote_average?: number;
  }
  
  export type Direction = 'book_to_movie' | 'movie_to_book' | 'book_to_book' | 'movie_to_movie';
  
  export interface RecommendRequest {
    selected_ids: string[];
    direction: Direction;
    weights: { genre: number; text: number; image: number };
  }
  
  export interface RecommendationResult {
    status: 'pending' | 'done' | 'error';
    movies?: string[];
    books?: string[];
    error?: string;
  }
  
  export interface HistoryEntry {
    id: string;
    task_id: string;
    selected_ids: string[];
    direction: Direction;
    weights: string;
    result: string;
    created_at: string;
  }
  
  export interface SavedRecommendation {
    id: string;
    from_type: 'book' | 'movie';
    from_id: string;
    to_type: 'book' | 'movie';
    to_id: string;
    saved_at: string;
  }