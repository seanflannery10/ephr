FROM golang:1.19.3-alpine AS builder
COPY .. /build/
WORKDIR /build
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o api ./cmd api

FROM scratch
COPY --from=builder /build/api /api
CMD ["/api"]
