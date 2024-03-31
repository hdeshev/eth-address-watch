FROM golang:1.22 as builder

RUN update-ca-certificates

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR "/src/eth-address-watch"

COPY go.mod *go.sum ./

RUN go mod download
COPY . .

RUN go build -v -o "../eth-address-watch"


# Second stage - `scratch` for production builds
FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

WORKDIR "/opt"

# Copy generated binary from previous image to this one - rename for readability
COPY --from=builder "/src/eth-address-watch" .

# Run the binary
CMD ["./eth-address-watch"]
