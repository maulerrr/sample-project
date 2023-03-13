FROM golang:alpine AS BUILDER
WORKDIR /app
COPY . .
RUN apk add build-base && go build -o forum api/main.go
FROM alpine:latest
WORKDIR /app
COPY --from=BUILDER /app .
CMD ["go run ./api/main.go"]
