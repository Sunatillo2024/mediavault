from pydantic import BaseModel
from typing import List, Optional
from datetime import datetime

class InstagramPost(BaseModel):
    id: str
    instagram_id: str
    caption: Optional[str] = None
    hashtags: List[str] = []
    media_type: str  # image, video, carousel
    media_urls: List[str] = []
    likes: int = 0
    comments: int = 0
    local_paths: List[str] = []
    created_at: Optional[str] = None