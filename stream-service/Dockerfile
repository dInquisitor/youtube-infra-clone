FROM golang:alpine
WORKDIR /stream-service/app
RUN go install github.com/mitranim/gow@latest

CMD [ "gow", "run", "." ]
    