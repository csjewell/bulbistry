# syntax=docker/dockerfile:1.4
FROM cgr.dev/chainguard/static:latest

LABEL org.opencontainers.image.title="Bulbistry"
LABEL org.opencontainers.image.description="Bulbistry minimal registry"
LABEL org.opencontainers.image.licenses=MIT
LABEL org.opencontainers.image.base.name=cgr.dev/chainguard/static:latest

COPY bulbistry /
COPY docker/blob/.* /blob/.* 
EXPOSE 8088
CMD ["/bulbistry"]


