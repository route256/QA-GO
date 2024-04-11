#!/bin/sh

GRPC_HOST="localhost:8082"

GRPC_METHOD="ozonmp.act_device_api.v1.ActDeviceApiService/DescribeDeviceV1"

payload=$(
cat <<EOF
  {
    "deviceId": 1
  }
EOF
)

grpcurl -plaintext -emit-defaults \
  -d "${payload}" ${GRPC_HOST} ${GRPC_METHOD}

GRPC_METHOD="ozonmp.act_device_api.v1.ActDeviceApiService/ListDevicesV1"
payload=$(
cat <<EOF
  {
    "page": 1,
    "perPage": 1
  }
EOF
)

grpcurl -plaintext -emit-defaults \
  -d "${payload}" ${GRPC_HOST} ${GRPC_METHOD}
