FROM gcr.io/distroless/static-debian12

COPY app /

ENTRYPOINT ["/app"]