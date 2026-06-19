# Multi-stage build for the cofiswarm-mode-cascade service.
FROM golang:1.25-alpine AS build
RUN apk add --no-cache git
WORKDIR /src
COPY go.mod go.sum* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -o /out/app ./cmd/cofiswarm-mode-cascade

FROM alpine:3.20
RUN adduser -D -u 10001 app
COPY --from=build /out/app /usr/local/bin/cofiswarm-mode-cascade
USER app
EXPOSE 8023
ENTRYPOINT ["/usr/local/bin/cofiswarm-mode-cascade"]
