FROM golang:alpine
WORKDIR /video-processor/app
RUN apk add ffmpeg
RUN go install github.com/mitranim/gow@latest

CMD [ "gow", "run", "." ]
    