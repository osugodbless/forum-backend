CREATE SCHEMA IF NOT EXISTS content;

CREATE TYPE content.reaction_type AS ENUM ('like', 'dislike');

-- Categories
CREATE TABLE IF NOT EXISTS content.post_categories (
    id uuid PRIMARY KEY,
    name varchar(100) UNIQUE NOT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Posts
CREATE TABLE IF NOT EXISTS content.posts (
    id uuid PRIMARY KEY,
    profile_id uuid NOT NULL, -- Logical reference to profile.profiles(id)
    title varchar(255) NOT NULL,
    content text NOT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Many-to-Many Map: Posts <-> Categories
CREATE TABLE IF NOT EXISTS content.post_category_map (
    post_id uuid NOT NULL REFERENCES content.posts(id) ON DELETE CASCADE,
    category_id uuid NOT NULL REFERENCES content.post_categories(id) ON DELETE CASCADE,
    PRIMARY KEY (post_id, category_id)
);

-- Comments (Supports nested replies within the domain)
CREATE TABLE IF NOT EXISTS content.comments (
    id uuid PRIMARY KEY,
    post_id uuid NOT NULL REFERENCES content.posts(id) ON DELETE CASCADE,
    profile_id uuid NOT NULL, -- Logical reference to profile.profiles(id)
    parent_comment_id uuid REFERENCES content.comments(id) ON DELETE CASCADE,
    content text NOT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Post Reactions
CREATE TABLE IF NOT EXISTS content.post_reactions (
    post_id uuid NOT NULL REFERENCES content.posts(id) ON DELETE CASCADE,
    profile_id uuid NOT NULL, -- Logical reference to profile.profiles(id)
    reaction_type content.reaction_type NOT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (post_id, profile_id)
);

-- Comment Reactions
CREATE TABLE IF NOT EXISTS content.comment_reactions (
    comment_id uuid NOT NULL REFERENCES content.comments(id) ON DELETE CASCADE,
    profile_id uuid NOT NULL, -- Logical reference to profile.profiles(id)
    reaction_type content.reaction_type NOT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (comment_id, profile_id)
);