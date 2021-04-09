FROM golang:1.16-alpine AS build_base

RUN apk add --no-cache git build-base

# Set the Current Working Directory inside the container
WORKDIR /tmp/commonpool

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

# Unit tests
RUN CGO_ENABLED=0 go test -v

# Build the Go app
RUN go build -o ./out/commonpool .

# Start fresh from a smaller image
FROM alpine:3.9
RUN apk add ca-certificates

COPY --from=build_base /tmp/commonpool/out/commonpool /app/commonpool
COPY public/ /app/

# Run the binary program produced by `go install`
CMD ["/app/commonpool"]