CREATE TABLE user_recommendation_history (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    task_id UUID NOT NULL,
    selected_ids TEXT NOT NULL,
    direction VARCHAR(20) NOT NULL,
    weights JSONB,
    result JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE user_saved_recommendations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    from_type VARCHAR(10) NOT NULL CHECK (from_type IN ('book','movie')),
    from_id VARCHAR(20) NOT NULL,
    to_type VARCHAR(10) NOT NULL CHECK (to_type IN ('book','movie')),
    to_id VARCHAR(20) NOT NULL,
    saved_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(user_id, from_type, from_id, to_type, to_id)
);