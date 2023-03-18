FROM golang:1.20 AS builder

RUN set -x \
	&& apt-get update \
	&& apt-get install -y --no-install-recommends \
		ca-certificates \
		git \
		libolm-dev

RUN echo 'nobody:x:65534:65534:nobody:/:' > /tmp/passwd \
	&& echo 'nobody:x:65534:' > /tmp/group

COPY go.mod go.sum /go/src/github.com/juli3nk/alertmanager-receiver-matrix/
WORKDIR /go/src/github.com/juli3nk/alertmanager-receiver-matrix

ENV GO111MODULE on
RUN go mod download

COPY . .

RUN go build -o /tmp/alertmanager-receiver-matrix


FROM debian:stable-slim

RUN set -x \
	&& apt-get update \
	&& apt-get install -y --no-install-recommends \
		ca-certificates \
		libolm-dev \
	&& apt-get -y autoremove \
	&& apt-get clean \
	&& rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

COPY --from=builder /tmp/alertmanager-receiver-matrix /usr/local/bin/alertmanager-receiver-matrix

USER nobody:nogroup

EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/alertmanager-receiver-matrix"]
