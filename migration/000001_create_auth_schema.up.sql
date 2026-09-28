CREATE TABLE IF NOT EXISTS auth.posts (
    id SERIAL PRIMARY KEY,
    profile_id INT NOT NULL REFERENCES auth.profiles(user_id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS auth.post_categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS auth.post_category_map (
    post_id INT NOT NULL REFERENCES auth.posts(id) ON DELETE CASCADE,
    category_id INT NOT NULL REFERENCES auth.post_categories(id) ON DELETE CASCADE,
    PRIMARY KEY (post_id, category_id)
);

CREATE TABLE IF NOT EXISTS auth.comments (
    id SERIAL PRIMARY KEY,
    post_id INT NOT NULL REFERENCES auth.posts(id) ON DELETE CASCADE,
    profile_id INT NOT NULL REFERENCES auth.profiles(user_id) ON DELETE CASCADE,
    parent_comment_id INT REFERENCES auth.comments(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TYPE auth.reaction_type AS ENUM ('like', 'dislike');

CREATE TABLE IF NOT EXISTS auth.post_reactions (
    post_id INT NOT NULL REFERENCES auth.posts(id) ON DELETE CASCADE,
    profile_id INT NOT NULL REFERENCES auth.profiles(user_id) ON DELETE CASCADE,
    reaction_type auth.reaction_type NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (post_id, profile_id)
);

CREATE TABLE IF NOT EXISTS auth.comment_reactions (
    comment_id INT NOT NULL REFERENCES auth.comments(id) ON DELETE CASCADE,
    profile_id INT NOT NULL REFERENCES auth.profiles(user_id) ON DELETE CASCADE,
    reaction_type auth.reaction_type NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (comment_id, profile_id)
);

CREATE INDEX IF NOT EXISTS idx_posts_profile_id ON auth.posts(profile_id);
CREATE INDEX IF NOT EXISTS idx_comments_post_id ON auth.comments(post_id);
CREATE INDEX IF NOT EXISTS idx_comments_profile_id ON auth.comments(profile_id);
CREATE INDEX IF NOT EXISTS idx_comments_parent_comment_id ON auth.comments(parent_comment_id);
CREATE INDEX IF NOT EXISTS idx_post_reactions_post_id ON auth.post_reactions(post_id);
CREATE INDEX IF NOT EXISTS idx_comment_reactions_comment_id ON auth.comment_reactions(comment_id);