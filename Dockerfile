FROM --platform=$BUILDPLATFORM golang:1.16 AS builder

LABEL org.opencontainers.image.source=https://github.com/mariamakbbh/slackbot_socketmode

ARG SKAFFOLD_GO_GCFLAGS
ARG TARGETOS
ARG TARGETARCH
ARG app_name=kafka

ENV GOOS=$TARGETOS
ENV GOARCH=$TARGETARCH
ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GO_BIN=/go/bin/app
ENV GRPC_HEALTH_PROBE_VERSION=v0.3.6

RUN apt-get update
RUN apt-get install -y build-essential git unzip curl
RUN curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip" 

RUN unzip awscliv2.zip
RUN ./aws/install
RUN aws --version
RUN wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
  chmod +x /bin/grpc_health_probe

WORKDIR /var/app

COPY . .

RUN make build

FROM gcr.io/distroless/static:nonroot

LABEL org.opencontainers.image.source=https://github.com/mariamakbbh/slackbot_socketmode

COPY --from=builder /go/bin/app /app
COPY --from=builder /bin/grpc_health_probe /bin/grpc_health_probe

CMD ["/app"]
