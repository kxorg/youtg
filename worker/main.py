from fastapi import FastAPI, HTTPException
from fastapi.responses import StreamingResponse
from pydantic import BaseModel
from pytube import YouTube
from io import BytesIO
import re
import logging

app = FastAPI()
logger = logging.getLogger(__name__)


class YouTubeURL(BaseModel):
    youtube_url: str


def is_valid_youtube_url(url: str) -> bool:
    youtube_regex = re.compile(
        r'(https?://)?(www\.)?(youtube|youtu|youtube-nocookie)\.(com|be)/'
        r'(watch\?v=|embed/|v/|.+\?v=)?([^&=%\?]{11})')
    return youtube_regex.match(url) is not None


@app.get("/")
async def index():
    return {"message": "ok"}


@app.post("/get_audio")
async def get_audio(item: YouTubeURL):
    if not is_valid_youtube_url(item.youtube_url):
        raise HTTPException(status_code=400, detail="Invalid YouTube URL")

    try:
        yt = YouTube(item.youtube_url)
        audio_stream = yt.streams.filter(only_audio=True).first()

        if audio_stream is None:
            raise HTTPException(status_code=404, detail="Audio stream not found")

        audio_buffer = BytesIO()
        audio_stream.stream_to_buffer(audio_buffer)
        audio_buffer.seek(0)
        return StreamingResponse(audio_buffer, media_type="audio/mp3")
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))
