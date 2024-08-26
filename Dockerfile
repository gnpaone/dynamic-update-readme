FROM golang:1.22 as builder

WORKDIR /app
COPY ./action /app

RUN go get -d -v

# Statically compile our app for use in a scratch or debian buster container
RUN CGO_ENABLED=0 go build -ldflags="-w -s" -v -o app .

# using multi-stage build
FROM debian:bookworm-slim
RUN set -x && apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    ca-certificates git bash && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/app /app

ENTRYPOINT ["/app"]
