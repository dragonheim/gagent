FROM golang:1.16-alpine3.12 as builder

WORKDIR /gagent
COPY . .
RUN apk add --no-cache zeromq-dev build-base git && \ 
    go build -o /gagent/bin/gagent gagent/main.go && \
    strip /gagent/bin/gagent


FROM alpine:3.12
LABEL Name="G'Agent"
LABEL Maintainer="jwells@dragonheim.net"
LABEL License="MIT License"
RUN apk add --no-cache zeromq; mkdir -p -m 0700 /etc/gagent
 
COPY --from=builder /gagent/examples/example_gagent.hcl /etc/gagent/gagent.hcl
COPY --from=builder /gagent/bin/gagent /usr/bin/

EXPOSE 35570/tcp
VOLUME /etc/gagent

CMD ["/usr/bin/gagent"]
