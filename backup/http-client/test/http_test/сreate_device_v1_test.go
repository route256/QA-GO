package http_test

import (
	"context"
	"crypto/rand"
	"math/big"
	"strconv"
	"testing"
	"time"

	apiClient "gitlab.ozon.dev/qa/classroom-4/act-device-api/test/http_test/client"
	"gitlab.ozon.dev/qa/classroom-4/act-device-api/test/http_test/models"

	"github.com/gobuffalo/envy"
	"github.com/stretchr/testify/assert"
)

func TestCreateDevice(t *testing.T) {

	var URL = envy.Get("BASE_URL", "http://127.0.0.1:8080")

	t.Run("Create device", func(t *testing.T) {
		// Arrange
		client := apiClient.NewHTTPClient(URL, 5, 1*time.Second)
		device := models.CreateDeviceRequest{
			Platform: "Ubuntu",
			UserID:   "701",
		}
		ctx := context.Background()

		// Act
		id, _, _ := client.CreateDevice(ctx, device)

		// Assert
		assert.GreaterOrEqual(t, id.DeviceID, 0)
	})

	t.Run("Create device and check description", func(t *testing.T) {

		// Arrange
		n, err := rand.Int(rand.Reader, big.NewInt(1000))
		if err != nil {
			t.Error("error:", err)
		}
		client := apiClient.NewHTTPClient(URL, 5, 1*time.Second)
		platform, userID := "Ubuntu", strconv.Itoa(int(n.Int64()))
		device := models.CreateDeviceRequest{
			Platform: platform,
			UserID:   userID,
		}
		ctx := context.Background()

		// Act
		id, _, _ := client.CreateDevice(ctx, device)
		description, _, _ := client.DescribeDevice(ctx, strconv.Itoa(id.DeviceID))

		// Assert
		assert.Equal(t, description.Value.ID, strconv.Itoa(id.DeviceID))
		assert.Equal(t, description.Value.Platform, platform)
		assert.Equal(t, description.Value.UserID, userID)
	})
}
