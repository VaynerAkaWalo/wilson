FROM golang:1.25.2-alpine AS build

WORKDIR /build

COPY go.mod .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o app cmd/main.go

FROM gcr.io/distroless/static-debian12

COPY --from=build build/app .
COPY --from=build build/migrations migrations

ENTRYPOINT ["/app"]