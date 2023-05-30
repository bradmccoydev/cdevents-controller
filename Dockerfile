FROM golang:1.20.4-alpine3.16 AS builder

ENV CGO_ENABLED=0
ARG REVISION
ARG VERSION
ARG COMMIT
ARG DATE
WORKDIR /workspace

COPY go.mod go.sum ./
RUN go mod download

COPY ./ ./

RUN go build -ldflags "-s -w \
    -X github.com/bradmccoydev/cdevents-controller/pkg/version.REVISION=${REVISION}" \
    -a -o /workspace/cdevents-controller cmd/cdevents-controller/*

RUN go build -ldflags "-s -w \
    -X github.com/bradmccoydev/cdevents-controller/pkg/version.REVISION=${REVISION}" \
    -a -o /workspace/cdeventscli cmd/cdeventscli/*

FROM gcr.io/distroless/static AS production

LABEL org.opencontainers.image.source="https://github.com/bradmccoydev/cdevents-controller" \
    org.opencontainers.image.url="https://avatars.githubusercontent.com/u/91484128?s=200&v=4" \
    org.opencontainers.image.title="CDEvents Github Controller" \
    org.opencontainers.image.vendor='bradmccoydev' \
    org.opencontainers.image.licenses='Apache-2.0' \
    org.opencontainers.image.description='CDEvents Controller'

WORKDIR /
COPY --from=builder /workspace/cdevents-controller .
COPY --from=builder /workspace/cdeventscli /usr/local/bin/cdeventscli
COPY ./ui ./ui
USER 65532:65532

ENTRYPOINT ["/cdevents-controller"]
