# Build Stage
FROM lacion/alpine-golang-buildimage:1.13 AS build-stage

LABEL app="build-go-hastily"
LABEL REPO="https://github.com/fhivemind/go-hastily"

ENV PROJPATH=/go/src/github.com/fhivemind/go-hastily

# Because of https://github.com/docker/docker/issues/14914
ENV PATH=$PATH:$GOROOT/bin:$GOPATH/bin

ADD . /go/src/github.com/fhivemind/go-hastily
WORKDIR /go/src/github.com/fhivemind/go-hastily

RUN make build-alpine

# Final Stage
FROM lacion/alpine-base-image:latest

ARG GIT_COMMIT
ARG VERSION
LABEL REPO="https://github.com/fhivemind/go-hastily"
LABEL GIT_COMMIT=$GIT_COMMIT
LABEL VERSION=$VERSION

# Because of https://github.com/docker/docker/issues/14914
ENV PATH=$PATH:/opt/go-hastily/bin

WORKDIR /opt/go-hastily/bin

COPY --from=build-stage /go/src/github.com/fhivemind/go-hastily/bin/go-hastily /opt/go-hastily/bin/
RUN chmod +x /opt/go-hastily/bin/go-hastily

# Create appuser
RUN adduser -D -g '' go-hastily
USER go-hastily

ENTRYPOINT ["/usr/bin/dumb-init", "--"]

CMD ["/opt/go-hastily/bin/go-hastily"]
