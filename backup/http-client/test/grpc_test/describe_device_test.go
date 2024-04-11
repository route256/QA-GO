package grpc_test

import (
	"context"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	act_device_api "gitlab.ozon.dev/qa/classroom-4/act-device-api/pkg/act-device-api/gitlab.ozon.dev/qa/classroom-4/act-device-api/pkg/act-device-api"

	"testing"

	"google.golang.org/grpc"
)

func TestDescribeDevice(t *testing.T) {
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

	t.Run("DescribeDevice existing", func(t *testing.T) {
		req := act_device_api.DescribeDeviceV1Request{
			DeviceId: 1,
		}

		res, err := actDeviceApiClient.DescribeDeviceV1(ctx, &req)
		require.NoError(t, err)
		require.NotNil(t, res)

		assert.EqualValues(t, res.Value.Id, 1)
		assert.Contains(t, "Android, Ios", res.Value.Platform)
		assert.NotEmpty(t, res.Value.UserId)
		assert.NotEmpty(t, res.Value.EnteredAt)
	})

	t.Run("DescribeDevice not existing", func(t *testing.T) {
		req := act_device_api.DescribeDeviceV1Request{
			DeviceId: 100,
		}

		_, err := actDeviceApiClient.DescribeDeviceV1(ctx, &req)
		assert.Equal(t, err.Error(), "rpc error: code = NotFound desc = device not found")
	})
}
