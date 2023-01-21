package http_test

import (
	"context"
	"net/url"
	"testing"
	"time"

	"github.com/gobuffalo/envy"
	"github.com/stretchr/testify/assert"
	apiClient "gitlab.ozon.dev/qa/classroom-4/act-device-api/test/http_test/client"
)

func TestListDevices(t *testing.T) {
	t.Run("Get devices", func(t *testing.T) {

		// Arrange
		var URL = envy.Get("BASE_URL", "http://127.0.0.1:8080")
		client := apiClient.NewHTTPClient(URL, 5, 1*time.Second)
		opts := url.Values{}
		opts.Add("page", "1")
		opts.Add("perPage", "100")
		ctx := context.Background()

		// Act
		items, _, _ := client.ListDevices(ctx, opts)

		// Assert
		assert.GreaterOrEqual(t, len(items.Items), 1)
	})
}
