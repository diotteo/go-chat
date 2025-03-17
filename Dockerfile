# syntax=docker/dockerfile:1
FROM debian:bookworm
COPY build/server /
EXPOSE 12345/udp
CMD ["/server"]
