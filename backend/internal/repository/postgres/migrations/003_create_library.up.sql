CREATE TABLE user_liked_books (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    book_id VARCHAR(20) NOT NULL,
    liked_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    PRIMARY KEY (user_id, book_id)
);

CREATE TABLE user_liked_movies (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    movie_id VARCHAR(20) NOT NULL,
    liked_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    PRIMARY KEY (user_id, movie_id)
);