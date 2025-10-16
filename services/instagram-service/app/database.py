import os
import asyncpg
import json
from app.models.post import InstagramPost
from typing import Optional

# Database connection pool
pool = None


async def get_connection_pool():
    global pool
    if pool is None:
        pool = await asyncpg.create_pool(
            host=os.getenv("DB_HOST", "localhost"),
            port=os.getenv("DB_PORT", "5432"),
            user=os.getenv("DB_USER", "admin"),
            password=os.getenv("DB_PASSWORD", "password123"),
            database=os.getenv("DB_NAME", "mediavault")
        )
    return pool


async def create_tables():
    """Create database tables if they don't exist"""
    pool = await get_connection_pool()

    create_table_sql = """
                       CREATE TABLE IF NOT EXISTS instagram_posts \
                       ( \
                           id \
                           VARCHAR \
                       ( \
                           255 \
                       ) PRIMARY KEY,
                           instagram_id VARCHAR \
                       ( \
                           255 \
                       ) NOT NULL,
                           caption TEXT,
                           hashtags JSONB,
                           media_type VARCHAR \
                       ( \
                           50 \
                       ),
                           media_urls JSONB,
                           likes INTEGER,
                           comments INTEGER,
                           local_paths JSONB,
                           created_at TIMESTAMP
                           ); \
                       """

    async with pool.acquire() as conn:
        await conn.execute(create_table_sql)


async def save_post(post: InstagramPost):
    """Save Instagram post to database"""
    pool = await get_connection_pool()

    insert_sql = """
                 INSERT INTO instagram_posts
                 (id, instagram_id, caption, hashtags, media_type, media_urls, likes, comments, local_paths, created_at)
                 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) \
                 """

    async with pool.acquire() as conn:
        await conn.execute(
            insert_sql,
            post.id,
            post.instagram_id,
            post.caption,
            json.dumps(post.hashtags),
            post.media_type,
            json.dumps(post.media_urls),
            post.likes,
            post.comments,
            json.dumps(post.local_paths),
            post.created_at
        )


async def get_post_by_id(post_id: str) -> Optional[InstagramPost]:
    """Get Instagram post by ID"""
    pool = await get_connection_pool()

    select_sql = """
                 SELECT id, \
                        instagram_id, \
                        caption, \
                        hashtags, \
                        media_type, \
                        media_urls, \
                        likes, \
                        comments, \
                        local_paths, \
                        created_at
                 FROM instagram_posts
                 WHERE id = $1 \
                 """

    async with pool.acquire() as conn:
        row = await conn.fetchrow(select_sql, post_id)

        if not row:
            return None

        return InstagramPost(
            id=row['id'],
            instagram_id=row['instagram_id'],
            caption=row['caption'],
            hashtags=json.loads(row['hashtags']),
            media_type=row['media_type'],
            media_urls=json.loads(row['media_urls']),
            likes=row['likes'],
            comments=row['comments'],
            local_paths=json.loads(row['local_paths']),
            created_at=row['created_at'].isoformat() if row['created_at'] else None
        )


async def get_post_by_instagram_id(instagram_id: str) -> Optional[InstagramPost]:
    """Get Instagram post by Instagram ID (shortcode)"""
    pool = await get_connection_pool()

    select_sql = """
                 SELECT id, \
                        instagram_id, \
                        caption, \
                        hashtags, \
                        media_type, \
                        media_urls, \
                        likes, \
                        comments, \
                        local_paths, \
                        created_at
                 FROM instagram_posts
                 WHERE instagram_id = $1 \
                 """

    async with pool.acquire() as conn:
        row = await conn.fetchrow(select_sql, instagram_id)

        if not row:
            return None

        return InstagramPost(
            id=row['id'],
            instagram_id=row['instagram_id'],
            caption=row['caption'],
            hashtags=json.loads(row['hashtags']),
            media_type=row['media_type'],
            media_urls=json.loads(row['media_urls']),
            likes=row['likes'],
            comments=row['comments'],
            local_paths=json.loads(row['local_paths']),
            created_at=row['created_at'].isoformat() if row['created_at'] else None
        )