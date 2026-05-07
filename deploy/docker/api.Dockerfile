FROM golang:1.26.2-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
	go build -o /out/api ./cmd/api

FROM alpine:3.22 AS runtime

WORKDIR /app

COPY --from=builder /out/api /app/api

EXPOSE 8080 8081

CMD ["/app/api"]
