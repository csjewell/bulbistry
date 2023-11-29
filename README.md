# Bulbistry

[![Go Report Card](https://goreportcard.com/badge/github.com/csjewell/bulbistry)](https://goreportcard.com/report/github.com/csjewell/bulbistry)
[![License](https://img.shields.io/github/license/csjewell/bulbistry)](./LICENSE.md)

The name comes from *Thomomys bulbivorus*, the scientific name for the Camas pocket gopher, and the fact that this implements an OCI-compliant container registry.

The intention is to write a registry for (semi-)private small-scale use to be wrapped in a container.
ONE container.
Not 10, like Harbor does.
Because of that, these requirements are set:

0) There is no **requirement** to use external services.
1) If you want HTTPS security, wrap a reverse proxy around this server.
2) The attached database is in SQLite format - no other databases will be supported until it DOES pass the spec.
3) Local storage only - if you want to hook up S3 or some other storage service, do it outside of this container.
4) This is implemented to the spec at https://github.com/opencontainers/distribution-spec/blob/main/spec.md, not to any specific program.
5) This is not a mirror. However, once it DOES pass the spec, we can talk about mirroring.
6) No support for vulnerability scanning will be entertained until it DOES pass the spec.

This is also my attempt to learn how to program in Go, as I've found that the best way for myself to learn is to read, and then try.
I've been programming in Perl for years, but I'm willing to try other languages.

Because of that, this is a work-in-progress, and is woefully incomplete even as yet.
