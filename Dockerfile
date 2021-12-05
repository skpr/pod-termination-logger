FROM alpine:3.15

RUN apk add --no-cache ca-certificates

COPY pod-termination-logger /usr/local/bin/
RUN chmod +x /usr/local/bin/pod-termination-logger

CMD ["/usr/local/bin/pod-termination-logger"]