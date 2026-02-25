FROM golang:1.24.3-alpine AS builder

RUN apk update && apk add --no-cache gcc musl-dev

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 \
    go build -installsuffix cgo \
    -ldflags="-s -w -extldflags '-static'" \
    -tags musl \
    -o jello-api ./cmd/main.go

FROM alpine:3.22.0 AS final

WORKDIR /app

COPY --from=builder /app/jello-api .

EXPOSE 3013

ENTRYPOINT ["./jello-api"]