# syntax=docker/dockerfile:1.4
FROM cgr.dev/chainguard/go:latest as build

WORKDIR /work
COPY . /work
RUN go build -o bulbistry cmd/bulbistry/main.go

FROM cgr.dev/chainguard/static:latest

LABEL org.opencontainers.image.title="Bulbistry"
LABEL org.opencontainers.image.description="Bulbistry minimal registry"
LABEL org.opencontainers.image.licenses=MIT
LABEL org.opencontainers.image.base.name=cgr.dev/chainguard/static:latest

COPY --from=build /work/bulbistry /bulbistry
EXPOSE 8088
CMD ["/bulbistry"]


