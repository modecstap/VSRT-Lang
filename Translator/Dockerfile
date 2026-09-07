FROM python:3.12-slim

WORKDIR /app/translator

COPY requirements.txt .

RUN pip install --no-cache-dir -r requirements.txt \
    && python -m nltk.downloader -d /usr/local/share/nltk_data wordnet

COPY . /app/translator

ENV PYTHONPATH=/app
ENV NLTK_DATA=/usr/local/share/nltk_data

EXPOSE 8080