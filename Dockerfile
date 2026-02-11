FROM golang:1.24 AS builder

ARG TARGETARCH

WORKDIR /build
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=${TARGETARCH} \
    go build -ldflags "-s -w" -o safecastbeat .

FROM ubuntu:24.04

COPY --from=builder /build/safecastbeat /usr/share/safecastbeat/safecastbeat
COPY safecastbeat.docker.yml /etc/safecastbeat/safecastbeat.yml
COPY safecastbeat.reference.yml /etc/safecastbeat/safecastbeat.reference.yml
COPY CA.pem /etc/safecastbeat/CA.pem
COPY fields.yml /usr/share/safecastbeat/fields.yml

ENTRYPOINT ["/usr/share/safecastbeat/safecastbeat"]
CMD ["-c", "/etc/safecastbeat/safecastbeat.yml", "-e"]
