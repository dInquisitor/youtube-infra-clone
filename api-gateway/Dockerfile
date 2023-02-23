FROM golang:alpine
WORKDIR /api-gateway/app
RUN go install github.com/mitranim/gow@latest

CMD [ "gow", "run", "." ]
    