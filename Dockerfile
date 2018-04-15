FROM alpine:3.6

ADD apiextensions-ca-helper_linux_amd64 /apiextensions-ca-helper

ENTRYPOINT ["/apiextensions-ca-helper"]
