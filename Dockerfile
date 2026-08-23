
# -----------------------------------------------------------------------------
# The is stage ONE of this multi-stage-build
# -----------------------------------------------------------------------------
FROM golang:1.26.5 AS builder

WORKDIR /app

# Copy file 'go.mod' from dev-environment to /app of this stage
COPY go.mod ./

# Copy file 'go.sum' from dev-environment to /app of this stage
COPY go.sum ./

# Download all dependencies for this project defined in 'go.mod' and 'go.sum'
RUN go mod download

# Copy ALL (other) files from dev-environment to /app of this stage
# Remember: All files exept the files/folders difined in '.dockerignore'
COPY . .

# 
RUN CGO_ENABLED=0 OOS=linux go build -o myapp ./cmd/server

# -----------------------------------------------------------------------------
# The is stage TWO of this multi-stage-build
# -----------------------------------------------------------------------------
# baseimage OHNE shell
# FROM gcr.io/distroless/static-debian12
# baseimage MIT shell
FROM alpine:3.20

# WORKDIR /app

COPY --from=builder /app/myapp /myapp

COPY --from=builder /app/certs/cert-key.pem /certs/cert-key.pem
COPY --from=builder /app/certs/cert.pem     /certs/cert.pem

# This app will use the following port ro listen on
EXPOSE 8080
EXPOSE 8443

# USER nonroot:nonroot

# CMD ["/myapp"]
ENTRYPOINT ["/myapp"]
