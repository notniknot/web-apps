FROM golang:1.24-alpine AS build
WORKDIR /src
COPY go.mod ./
COPY cmd ./cmd
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/web ./cmd/web

FROM alpine:3.22
RUN adduser -D -u 10001 app
USER 10001
COPY --from=build /out/web /web
EXPOSE 8080
ENTRYPOINT ["/web"]

