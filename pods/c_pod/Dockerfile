FROM gcc:latest

WORKDIR /app

ADD . /app

LABEL type="XML"

RUN gcc -o server server.c

EXPOSE 3032

CMD ["./server"]