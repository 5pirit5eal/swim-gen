### Build Go server binary
FROM golang:1.26@sha256:24a707af49702fa44f02785475d02fb752f6982d3e9c0afb7696e1f0204e7ce9 AS builder

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . ./

# CGO_ENABLED=0 disables cgo, which allows to build a static binary.
RUN CGO_ENABLED=0 GOOS=linux go build -mod=readonly -v -o service ./main.go

# Prepare final image
# static-debian12 is a distroless image including ca-certificates.
FROM gcr.io/distroless/static-debian12@sha256:6447365a6337c3732f412d1b74357b30a633831955b2bc45552b0086be907687

# Copy Go binary
COPY --from=builder /build/service /bin

CMD ["/bin/service"]
