FROM golang:1.19.3-alpine AS builder

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o api ./cmd api

FROM scratch
COPY --from=builder /api /api
CMD ["/api"]
