FROM golang:1.9.6 as builder

RUN mkdir /go/src/api_exporter

COPY [".", "/go/src/api_exporter"]

ENV https_proxy="http://Guys:EscapeIran@178.128.164.255:7777" \
    http_proxy="http://Guys:EscapeIran@178.128.164.255:7777"
    

RUN cd /go/src/api_exporter/ && go get && CGO_ENABLED=0 go build -o /api_exporter

FROM busybox:1.28.4 as api_exporter

COPY --from=builder /api_exporter /

RUN chmod +x /api_exporter

EXPOSE 8283

ENTRYPOINT ["/api_exporter"]
