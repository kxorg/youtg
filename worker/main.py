from fastapi import FastAPI
from pydantic import BaseModel

app = FastAPI()


class TextItem(BaseModel):
    text: str


@app.post("/echo")
async def echo_text(item: TextItem):
    return {"text": item.text}
