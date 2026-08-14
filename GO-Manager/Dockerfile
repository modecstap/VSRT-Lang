# Build stage
FROM golang:1.26-alpine AS builder
WORKDIR /app

# install git for go modules if needed
RUN apk add --no-cache git

# copy go modules files
COPY go.mod go.sum ./
RUN go mod download

# copy source
COPY . .

# build the binary (builds package in ./cmd)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /usr/local/bin/gomanager ./cmd

# Final stage: minimal runtime image
FROM scratch
COPY --from=builder /usr/local/bin/gomanager /usr/local/bin/gomanager

EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/gomanager"]
