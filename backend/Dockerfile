### Build Go server binary
FROM golang:1.25@sha256:d7ae6dbb143c4eb97788a53a2bfa879cdde91f793c1d44d9f3bc7cb0df2e4e5f AS builder

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
