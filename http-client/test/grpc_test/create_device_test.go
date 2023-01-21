package grpc_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	act_device_api "gitlab.ozon.dev/qa/classroom-4/act-device-api/pkg/act-device-api/gitlab.ozon.dev/qa/classroom-4/act-device-api/pkg/act-device-api"
	"google.golang.org/grpc"
)

func TestCreateDevice(t *testing.T) {
	host := "localhost:8082"
	ctx := context.Background()
	conn, err := grpc.Dial(host, grpc.WithInsecure())
	if err != nil {
		t.Fatalf("grpc.Dial() err: %v", err)
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			t.Logf("conn.Close err: %v", err)
		}
	}(conn)

	actDeviceApiClient := act_device_api.NewActDeviceApiServiceClient(conn)

	t.Run("CreateDevice valid", func(t *testing.T) {
		req := act_device_api.CreateDeviceV1Request{
			Platform: "Android",
			UserId:   123000,
		}

		res, err := actDeviceApiClient.CreateDeviceV1(ctx, &req)
		require.NoError(t, err)
		require.NotNil(t, res)

		assert.GreaterOrEqual(t, res.DeviceId, uint64(1))
	})

	t.Run("CreateDevice invalid", func(t *testing.T) {
		req := act_device_api.CreateDeviceV1Request{
			Platform: "",
			UserId:   666,
		}

		_, err := actDeviceApiClient.CreateDeviceV1(ctx, &req)
		assert.Equal(t, err.Error(), "rpc error: code = InvalidArgument desc = invalid CreateDeviceV1Request.Platform: value length must be at least 1 runes")
	})
}
