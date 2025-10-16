from fastapi import FastAPI
from app.api.routes import router
from app.database import create_tables
import os

app = FastAPI(
    title="Instagram Service",
    description="Microservice for Instagram content extraction",
    version="1.0.0"
)

# Include routes
app.include_router(router)

@app.on_event("startup")
async def startup_event():
    await create_tables()

@app.get("/")
async def root():
    return {"message": "Instagram Service is running"}

if __name__ == "__main__":
    import uvicorn
    port = int(os.getenv("SERVICE_PORT", "8082"))
    uvicorn.run(app, host="0.0.0.0", port=port)