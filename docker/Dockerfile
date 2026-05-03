# syntax=docker/dockerfile:1.20
FROM cgr.dev/chainguard/static:latest
ARG TARGETPLATFORM

COPY bulbistry /
COPY docker/blob/.* /blob/.* 
EXPOSE 8088
CMD ["/bulbistry"]


