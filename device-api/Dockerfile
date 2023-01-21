# Builder

ARG GITHUB_PATH=gitlab.ozon.dev/qa/classroom-4/act-device-api

FROM golang:1.19.1-alpine AS builder
RUN apk add --update make curl git protoc protobuf protobuf-dev
COPY . /home/${GITHUB_PATH}
WORKDIR /home/${GITHUB_PATH}
RUN make deps-go && make build-go

# gRPC Server

FROM alpine:latest as server
LABEL org.opencontainers.image.source https://${GITHUB_PATH}
RUN apk --no-cache add ca-certificates
RUN apk --no-cache add curl
WORKDIR /root/

COPY --from=builder /home/${GITHUB_PATH}/bin/grpc-server .
COPY --from=builder /home/${GITHUB_PATH}/config.yml .
COPY --from=builder /home/${GITHUB_PATH}/migrations/ ./migrations

RUN chown root:root grpc-server

EXPOSE 50051
EXPOSE 8080
EXPOSE 9100

CMD ["./grpc-server"]
