FROM python:3.11-alpine

WORKDIR /app
RUN apk add --no-cache build-base libffi-dev openssl-dev
COPY requirements.txt .
RUN python -m venv /venv \
    && /venv/bin/pip install --no-cache-dir -r requirements.txt
COPY . .
ENV PATH="/venv/bin:$PATH"
EXPOSE 8080
#CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "8080"]
CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "8080", "--reload"]
