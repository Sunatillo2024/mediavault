from fastapi import APIRouter, HTTPException
from pydantic import BaseModel
from typing import Optional, List
from app.services.instagram_service import InstagramService
import uuid

router = APIRouter(prefix="/api", tags=["instagram"])

instagram_service = InstagramService()

class ExtractRequest(BaseModel):
    url: str
    download_media: bool = True

class InstagramPostResponse(BaseModel):
    id: str
    instagram_id: str
    caption: Optional[str]
    hashtags: List[str]
    media_type: str
    media_urls: List[str]
    likes: int
    comments: int
    local_paths: List[str]
    created_at: str

@router.get("/health")
async def health_check():
    return {"status": "ok", "service": "instagram"}

@router.post("/extract", response_model=InstagramPostResponse)
async def extract_post(request: ExtractRequest):
    try:
        post = await instagram_service.extract_post(
            request.url,
            request.download_media
        )
        return post
    except ValueError as e:
        raise HTTPException(status_code=400, detail=str(e))
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Extraction failed: {str(e)}")

@router.get("/post/{post_id}", response_model=InstagramPostResponse)
async def get_post(post_id: str):
    try:
        # Validate UUID
        uuid.UUID(post_id)
        post = await instagram_service.get_post(post_id)
        if not post:
            raise HTTPException(status_code=404, detail="Post not found")
        return post
    except ValueError:
        raise HTTPException(status_code=400, detail="Invalid post ID")
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

@router.get("/media/{post_id}")
async def get_media_paths(post_id: str):
    try:
        uuid.UUID(post_id)
        post = await instagram_service.get_post(post_id)
        if not post:
            raise HTTPException(status_code=404, detail="Post not found")
        return {"media_paths": post.local_paths}
    except ValueError:
        raise HTTPException(status_code=400, detail="Invalid post ID")
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))