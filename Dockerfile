FROM golang:alpine

ADD ./server /server
WORKDIR /server

RUN apk add --no-cache git

RUN go build

CMD ["./server"]