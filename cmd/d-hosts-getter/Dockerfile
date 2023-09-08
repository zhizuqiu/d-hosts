FROM docker.io/alpine:latest

COPY dist/amd64/d-hosts-getter /

RUN chmod +x /d-hosts-getter

ENTRYPOINT ["/d-hosts-getter"]