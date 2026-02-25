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
    -o jello-api ./main.go


FROM alpine:3.22.0 AS final

WORKDIR /app

# Install curl
RUN apk add --no-cache curl

# Install Doppler CLI
RUN curl -Ls https://cli.doppler.com/install.sh | sh

COPY --from=builder /app/jello-api .

EXPOSE 3013

# Use doppler run to inject secrets
ENTRYPOINT ["doppler", "run", "--", "./jello-api"]