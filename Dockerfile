### Build Go server binary
FROM golang:1.24 AS builder

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . ./

# CGO_ENABLED=0 disables cgo, which allows to build a static binary.
RUN CGO_ENABLED=0 GOOS=linux go build -mod=readonly -v -o service ./cmd/swim-rag/swim-rag.go

# Prepare final image
# static-debian11 is a distroless image (approx. 2MiB) including ca-certificates.
FROM gcr.io/distroless/static-debian12

# Copy Go binary
COPY --from=builder /build/service /bin

CMD ["/bin/service"]