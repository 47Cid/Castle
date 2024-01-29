FROM python:3.7-slim

WORKDIR /app

ADD . /app

LABEL type="SQL"

LABEL weight="1"

COPY requirements.txt requirements.txt

RUN pip install --no-cache-dir -r requirements.txt

EXPOSE 3031

CMD ["python", "app.py"]