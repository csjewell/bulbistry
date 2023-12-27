# syntax=docker/dockerfile:1.4
FROM cgr.dev/chainguard/go:latest as build

WORKDIR /work
COPY . /work
RUN go build -o bulbistry cmd/bulbistry/main.go

FROM cgr.dev/chainguard/static:latest

LABEL org.opencontainers.image.source=https://github.com/csjewell/bulbistry
LABEL org.opencontainers.image.description="Bulbistry minimal registry"
LABEL org.opencontainers.image.licenses=MIT

COPY --from=build /work/bulbistry /bulbistry
EXPOSE 8088
CMD ["/bulbistry"]


