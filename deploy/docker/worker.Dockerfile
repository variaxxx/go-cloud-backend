FROM golang:1.26.2-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
	go build -o /out/worker ./cmd/worker

FROM alpine:3.22 AS runtime

WORKDIR /app

RUN apk add poppler-utils
COPY --from=builder /out/worker /app/worker

CMD ["/app/worker"]
