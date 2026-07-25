FROM python:3.12-slim

WORKDIR /app/translator

COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

COPY . /app/translator
ENV PYTHONPATH=/app

EXPOSE 8080