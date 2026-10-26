CREATE SCHEMA IF NOT EXISTS content;

CREATE TYPE content.reaction_type AS ENUM ('like', 'dislike');

-- Categories
CREATE TABLE IF NOT EXISTS content.post_categories (
    id uuid PRIMARY KEY,
    name text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT uq_post_categories_name UNIQUE (name)
);

-- Posts
CREATE TABLE IF NOT EXISTS content.posts (
    id uuid PRIMARY KEY,
    profile_id uuid NOT NULL, -- Logical reference to profile.profiles(id)
    title text,
    content text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT chk_posts_title_length CHECK (title IS NULL OR char_length(trim(title)) BETWEEN 3 AND 255),
    CONSTRAINT chk_posts_content_not_empty CHECK (char_length(trim(content)) BETWEEN 7 AND 3000)
);

-- Lookup indexes for author history and reverse chronological feeds
CREATE INDEX IF NOT EXISTS idx_posts_profile_id ON content.posts (profile_id);
CREATE INDEX IF NOT EXISTS idx_posts_created_at_desc ON content.posts (created_at DESC);


-- Many-to-Many Map: Posts <-> Categories
CREATE TABLE IF NOT EXISTS content.post_category_map (
    post_id uuid NOT NULL,
    category_id uuid NOT NULL,
    PRIMARY KEY (post_id, category_id),

    CONSTRAINT fk_post_category_map_post FOREIGN KEY (post_id) REFERENCES content.posts(id) ON DELETE CASCADE,
    CONSTRAINT fk_post_category_map_category FOREIGN KEY (category_id) REFERENCES content.post_categories(id) ON DELETE CASCADE
);

-- Index for finding posts under a category
CREATE INDEX IF NOT EXISTS idx_post_category_map_category_id ON content.post_category_map (category_id);

-- Comments (Supports nested replies)
CREATE TABLE IF NOT EXISTS content.comments (
    id uuid PRIMARY KEY,
    post_id uuid NOT NULL,
    profile_id uuid NOT NULL, -- Logical reference to profile.profiles(id)
    parent_comment_id uuid,
    content text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_comments_post FOREIGN KEY (post_id) REFERENCES content.posts(id) ON DELETE CASCADE,
    CONSTRAINT fk_comments_parent FOREIGN KEY (parent_comment_id) REFERENCES content.comments(id) ON DELETE CASCADE,
    CONSTRAINT chk_comments_content_length CHECK (char_length(trim(content)) BETWEEN 1 AND 3000)
);

-- Indexes for post comments, threaded comment trees, and user activity
CREATE INDEX IF NOT EXISTS idx_comments_post_id ON content.comments (post_id);
CREATE INDEX IF NOT EXISTS idx_comments_parent_id ON content.comments (parent_comment_id);
CREATE INDEX IF NOT EXISTS idx_comments_profile_id ON content.comments (profile_id);

-- Post Reactions
CREATE TABLE IF NOT EXISTS content.post_reactions (
    post_id uuid NOT NULL,
    profile_id uuid NOT NULL, -- Logical reference to profile.profiles(id)
    reaction_type content.reaction_type NOT NULL,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (post_id, profile_id),

    CONSTRAINT fk_post_reactions_post FOREIGN KEY (post_id) REFERENCES content.posts(id) ON DELETE CASCADE
);

-- Comment Reactions
CREATE TABLE IF NOT EXISTS content.comment_reactions (
    comment_id uuid NOT NULL,
    profile_id uuid NOT NULL, -- Logical reference to profile.profiles(id)
    reaction_type content.reaction_type NOT NULL,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (comment_id, profile_id),

    CONSTRAINT fk_comment_reactions_comment FOREIGN KEY (comment_id) REFERENCES content.comments(id) ON DELETE CASCADE
);