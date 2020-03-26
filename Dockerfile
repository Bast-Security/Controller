FROM golang:1.13

EXPOSE 8080/tcp

ENV BAST_DB_DB bast
ENV BAST_DB_USER bast
ENV BAST_DB_PASS bast
ENV BAST_CERT ./pki/bast.crt
ENV BAST_KEY ./pki/bast.key

WORKDIR /go/src/bast
COPY . .

RUN go build

CMD ["./controller"]
