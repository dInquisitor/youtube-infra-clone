FROM golang:alpine
WORKDIR /api-server/app
RUN go install github.com/mitranim/gow@latest
# RUN go get gopkg.in/gorp.v1
# RUN go get github.com/codegangsta/gin
# RUN go get github.com/lib/pq
# RUN go get github.com/gin-gonic/gin
CMD [ "gow", "run", "." ]
    