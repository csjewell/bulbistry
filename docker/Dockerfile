# syntax=docker/dockerfile:1.4
FROM cgr.dev/chainguard/go:latest as build

WORKDIR /work
COPY . /work
RUN go build -o bulbistry .

FROM cgr.dev/chainguard/static:latest

COPY --from=build /work/bulbistry /bulbistry
EXPOSE 8088
CMD ["/bulbistry"]


