FROM golang:1.13

EXPOSE 8080

WORKDIR /go/src/bast
COPY . .

RUN go build

CMD ["./controller"]
