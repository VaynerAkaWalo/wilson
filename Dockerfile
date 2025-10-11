FROM gcr.io/distroless/static-debian12

COPY build/app /app

ENTRYPOINT ["/app"]