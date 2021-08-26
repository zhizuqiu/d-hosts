FROM docker.io/alpine:latest

COPY d-hosts-getter /

RUN chmod +x /d-hosts-getter

ENTRYPOINT ["/d-hosts-getter"]