FROM golang:1.13

EXPOSE 8080/tcp
EXPOSE 5353/udp

WORKDIR /go/src/bast
COPY . .

RUN go build

CMD ["./controller"]
