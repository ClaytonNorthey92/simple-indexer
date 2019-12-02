FROM golang:1.12-alpine as build

RUN apk add git

WORKDIR /build

COPY . .

ENV CGO_ENABLED=0

RUN go test -v ./...

RUN go build -o /build/server ./...

FROM scratch

# needed to make outgoing requests, doesn't come in scratch
COPY --from=build /etc/ssl/certs /etc/ssl/certs

WORKDIR /app

COPY --from=build /build/server /app

CMD ["/app/server"]