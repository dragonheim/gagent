# FROM golang:1.17-alpine3.14 as builder
FROM dragonheim/golang:1.17.0 as builder
ARG SEMVER

WORKDIR /gagent
COPY . .

ARG GOOS=${GOOS:-linux}
# ARG GOARCH=${GOARCH:-amd64}
ARG CGO_ENABLED=1

RUN apk add --no-cache zeromq-dev build-base git
RUN go build -o /gagent/bin/gagent -ldflags "-X main.semVER=${SEMVER}" gagent/main.go
RUN strip /gagent/bin/gagent

FROM alpine:3.14
ARG vernum
LABEL Name="G'Agent"
LABEL Maintainer="jwells@dragonheim.net"
LABEL License="MIT License"
LABEL version="${SEMVER}"
RUN apk add --no-cache zeromq && mkdir -p -m 0700 /etc/gagent
 
COPY --from=builder /gagent/examples/gagent.hcl /etc/gagent/gagent.hcl
COPY --from=builder /gagent/bin/gagent /usr/bin/

#      Router    Client    Worker
EXPOSE 35570/tcp 35571/tcp 35572/tcp
VOLUME /etc/gagent

CMD ["/usr/bin/gagent"]

HEALTHCHECK NONE
