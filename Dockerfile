FROM golang:alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o repimage .

FROM alpine
COPY --from=builder /app/repimage /repimage
COPY ./certs /certs

# Bake the default registry mappings into the image so the webhook can load it at runtime.
# This avoids relying on the container working directory ("./config/registries.json").
RUN mkdir -p /etc/repimage
COPY ./config/registries.json /etc/repimage/registries.json
