import os
import re
import uuid
import instaloader
from typing import List, Optional
import requests
from app.models.post import InstagramPost
from app.database import save_post, get_post_by_id, get_post_by_instagram_id
import aiohttp
import asyncio
from pathlib import Path


class InstagramService:
    def __init__(self):
        self.loader = instaloader.Instaloader()
        self.storage_path = os.getenv("STORAGE_PATH", "./storage")

    async def extract_post(self, url: str, download_media: bool = True) -> InstagramPost:
        """Extract Instagram post from URL"""

        # Extract shortcode from URL
        shortcode = self._extract_shortcode(url)
        if not shortcode:
            raise ValueError("Invalid Instagram URL")

        # Check if post already exists
        existing_post = await get_post_by_instagram_id(shortcode)
        if existing_post:
            return existing_post

        try:
            # Get post using instaloader
            post = instaloader.Post.from_shortcode(self.loader.context, shortcode)

            # Extract hashtags from caption
            hashtags = self._extract_hashtags(post.caption) if post.caption else []

            # Get media URLs
            media_urls = await self._get_media_urls(post)

            # Create post object
            instagram_post = InstagramPost(
                id=str(uuid.uuid4()),
                instagram_id=shortcode,
                caption=post.caption,
                hashtags=hashtags,
                media_type=self._get_media_type(post),
                media_urls=media_urls,
                likes=post.likes,
                comments=post.comments,
                local_paths=[],
                created_at=post.date_utc.isoformat() if post.date_utc else None
            )

            # Download media if requested
            if download_media and media_urls:
                local_paths = await self._download_media(media_urls, instagram_post.id)
                instagram_post.local_paths = local_paths

            # Save to database
            await save_post(instagram_post)

            return instagram_post

        except Exception as e:
            raise ValueError(f"Failed to extract post: {str(e)}")

    def _extract_shortcode(self, url: str) -> Optional[str]:
        """Extract shortcode from Instagram URL"""
        patterns = [
            r"instagram\.com/p/([^/?]+)",
            r"instagram\.com/reel/([^/?]+)",
            r"instagram\.com/stories/[^/]+/([^/?]+)"
        ]

        for pattern in patterns:
            match = re.search(pattern, url)
            if match:
                return match.group(1)
        return None

    def _extract_hashtags(self, caption: str) -> List[str]:
        """Extract hashtags from caption"""
        return re.findall(r"#(\w+)", caption)

    def _get_media_type(self, post: instaloader.Post) -> str:
        """Determine media type"""
        if post.is_video:
            return "video"
        elif post.mediacount > 1:
            return "carousel"
        else:
            return "image"

    async def _get_media_urls(self, post: instaloader.Post) -> List[str]:
        """Get all media URLs from post"""
        media_urls = []

        if post.mediacount == 1:
            # Single media post
            media_urls.append(post.url)
        else:
            # Carousel post
            for node in post.get_sidecar_nodes():
                media_urls.append(node.display_url)

        return media_urls

    async def _download_media(self, media_urls: List[str], post_id: str) -> List[str]:
        """Download media files to local storage"""
        local_paths = []
        media_dir = Path(self.storage_path) / "instagram" / "media" / post_id
        media_dir.mkdir(parents=True, exist_ok=True)

        async with aiohttp.ClientSession() as session:
            for i, url in enumerate(media_urls):
                try:
                    file_extension = self._get_file_extension(url)
                    filename = f"media_{i + 1}{file_extension}"
                    filepath = media_dir / filename

                    async with session.get(url) as response:
                        if response.status == 200:
                            content = await response.read()
                            with open(filepath, 'wb') as f:
                                f.write(content)
                            local_paths.append(str(filepath))
                except Exception as e:
                    print(f"Failed to download media {url}: {str(e)}")

        return local_paths

    def _get_file_extension(self, url: str) -> str:
        """Get file extension from URL"""
        if '.jpg' in url or '.jpeg' in url:
            return '.jpg'
        elif '.png' in url:
            return '.png'
        elif '.mp4' in url:
            return '.mp4'
        else:
            return '.jpg'  # default

    async def get_post(self, post_id: str) -> Optional[InstagramPost]:
        """Get post by ID"""
        return await get_post_by_id(post_id)