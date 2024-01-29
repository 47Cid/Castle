FROM gcc:latest

WORKDIR /app

ADD . /app

LABEL type="DOS"

LABEL weight=1

RUN gcc -o server server.c

EXPOSE 3032

CMD ["./server"]