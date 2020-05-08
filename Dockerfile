# build the binary
FROM golang:alpine AS builder

WORKDIR /go/src/app
COPY . .
RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -v -o app

# build a small image that runs the binary
FROM scratch

COPY --from=builder /go/src/app/app .

CMD ["./app"]