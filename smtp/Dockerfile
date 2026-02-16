FROM golang:1.21-alpine AS builder

WORKDIR /app

RUN go mod init smtp-manager && \
    go get github.com/gin-gonic/gin@latest && \
    go get github.com/lib/pq@latest && \
    go get golang.org/x/crypto@latest

COPY main.go .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o smtp-manager main.go

FROM alpine:3.18

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/smtp-manager .

EXPOSE 8080

CMD ["./smtp-manager"]
