package http_test

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/gobuffalo/envy"
	"github.com/stretchr/testify/assert"
	apiClient "gitlab.ozon.dev/qa/classroom-4/act-device-api/test/http_test/client"
	"gitlab.ozon.dev/qa/classroom-4/act-device-api/test/http_test/models"
)

func TestDeleteDevice(t *testing.T) {

	var URL = envy.Get("BASE_URL", "http://127.0.0.1:8080")

	t.Run("Delete device", func(t *testing.T) {
		// Arrange
		client := apiClient.NewHTTPClient(URL, 5, 1*time.Second)
		device := models.CreateDeviceRequest{
			Platform: "Ubuntu",
			UserID:   "701",
		}
		ctx := context.Background()

		// Act
		id, _, _ := client.CreateDevice(ctx, device)
		deletedDevice, _, _ := client.RemoveDevice(ctx, strconv.Itoa(id.DeviceID))

		// Assert
		assert.Equal(t, deletedDevice.Found, true, "Device deleted")
	})
}
