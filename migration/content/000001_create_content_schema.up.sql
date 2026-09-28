CREATE SCHEMA IF NOT EXISTS content;

CREATE TYPE content.reaction_type AS ENUM ('like', 'dislike');

-- Categories
CREATE TABLE IF NOT EXISTS content.post_categories (
    id UUID PRIMARY KEY DEFAULT,
    name VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Posts
CREATE TABLE IF NOT EXISTS content.posts (
    id UUID PRIMARY KEY DEFAULT,
    profile_id UUID NOT NULL, -- Logical reference to profile.profiles(id)
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Many-to-Many Map: Posts <-> Categories
CREATE TABLE IF NOT EXISTS content.post_category_map (
    post_id UUID NOT NULL REFERENCES content.posts(id) ON DELETE CASCADE,
    category_id UUID NOT NULL REFERENCES content.post_categories(id) ON DELETE CASCADE,
    PRIMARY KEY (post_id, category_id)
);

-- Comments (Supports nested replies within the domain)
CREATE TABLE IF NOT EXISTS content.comments (
    id UUID PRIMARY KEY DEFAULT,
    post_id UUID NOT NULL REFERENCES content.posts(id) ON DELETE CASCADE,
    profile_id UUID NOT NULL, -- Logical reference to profile.profiles(id)
    parent_comment_id UUID REFERENCES content.comments(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Post Reactions
CREATE TABLE IF NOT EXISTS content.post_reactions (
    post_id UUID NOT NULL REFERENCES content.posts(id) ON DELETE CASCADE,
    profile_id UUID NOT NULL, -- Logical reference to profile.profiles(id)
    reaction_type content.reaction_type NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (post_id, profile_id)
);

-- Comment Reactions
CREATE TABLE IF NOT EXISTS content.comment_reactions (
    comment_id UUID NOT NULL REFERENCES content.comments(id) ON DELETE CASCADE,
    profile_id UUID NOT NULL, -- Logical reference to profile.profiles(id)
    reaction_type content.reaction_type NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (comment_id, profile_id)
);