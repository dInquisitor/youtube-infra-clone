FROM golang:alpine
WORKDIR /upload-service/app
RUN go install github.com/mitranim/gow@latest

CMD [ "gow", "run", "." ]
