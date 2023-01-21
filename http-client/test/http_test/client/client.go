package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	models2 "gitlab.ozon.dev/qa/classroom-4/act-device-api/test/http_test/models"

	"github.com/hashicorp/go-retryablehttp"
	"gitlab.ozon.dev/qa/classroom-4/act-device-api/internal/pkg/logger"
)

// Client is a client for the Act Device API
type Client interface {
	Do(req *http.Request) (*http.Response, error)
	CreateDevice(ctx context.Context, body models2.CreateDeviceRequest) (models2.CreateDeviceResponse, *http.Response, error)
	ListDevices(ctx context.Context, opts url.Values) (models2.ListItemsResponse, *http.Response, error)
	DescribeDevice(ctx context.Context, deviceID string) (models2.DescribeDeviceResponse, *http.Response, error)
	RemoveDevice(ctx context.Context, deviceID string) (models2.RemovedDevice, *http.Response, error)
}

type client struct {
	client   *retryablehttp.Client
	BasePath string
}

// NewHTTPClient creates a new HTTP client.
func NewHTTPClient(basePath string, retryMax int, timeout time.Duration) Client {
	c := &retryablehttp.Client{
		HTTPClient:      &http.Client{Timeout: timeout},
		RetryMax:        retryMax,
		RetryWaitMin:    1 * time.Second,
		RetryWaitMax:    10 * time.Second,
		CheckRetry:      retryablehttp.DefaultRetryPolicy,
		Backoff:         retryablehttp.DefaultBackoff,
		RequestLogHook:  requestHook,
		ResponseLogHook: responseHook,
	}

	client := &client{client: c, BasePath: basePath}
	return client
}

func requestHook(_ retryablehttp.Logger, req *http.Request, retry int) {
	dump, err := httputil.DumpRequest(req, true) // better way
	if err != nil {
		logger.ErrorKV(req.Context(), "can't dump request")
	}

	logger.InfoKV(
		req.Context(),
		fmt.Sprintf("Retry request %d", retry),
		"request", string(dump),
		"url", req.URL.String(),
	)
}

func responseHook(_ retryablehttp.Logger, res *http.Response) {
	dump, err := httputil.DumpResponse(res, true) // better way
	if err != nil {
		logger.ErrorKV(res.Request.Context(), "can't dump response")
	}

	logger.InfoKV(
		res.Request.Context(),
		"Responded",
		"response", dump,
		"url", res.Request.URL.String(),
		"status_code", res.StatusCode,
	)
}

func (c *client) Do(request *http.Request) (*http.Response, error) {
	req, err := retryablehttp.FromRequest(request)
	if err != nil {
		return nil, err
	}

	return c.client.Do(req)
}

func (c *client) CreateDevice(ctx context.Context, body models2.CreateDeviceRequest) (models2.CreateDeviceResponse, *http.Response, error) {
	var localResponse models2.CreateDeviceResponse

	b := new(bytes.Buffer)
	err := json.NewEncoder(b).Encode(body)
	if err != nil {
		return localResponse, nil, err
	}

	req, err := http.NewRequest(http.MethodPost, c.BasePath+"/api/v1/devices", b)
	if err != nil {
		return localResponse, nil, err
	}
	res, err := c.Do(req)
	if err != nil {
		return localResponse, res, err
	}

	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			logger.ErrorKV(ctx, "Error on Body reading", err)
		}
	}(res.Body)

	if res.StatusCode != http.StatusOK {
		logger.ErrorKV(ctx, "Bad status code", res.StatusCode)
		return localResponse, res, err
	}
	data, _ := ioutil.ReadAll(res.Body)
	device := new(models2.CreateDeviceResponse)
	err = json.Unmarshal(data, &device)
	if err != nil {
		return localResponse, res, err
	}
	return *device, res, nil
}

func (c *client) ListDevices(ctx context.Context, opts url.Values) (models2.ListItemsResponse, *http.Response, error) {
	var localResponse models2.ListItemsResponse

	apiURL, err := url.Parse(c.BasePath + "/api/v1/devices")
	if err != nil {
		return localResponse, nil, err
	}

	query := apiURL.Query()
	for k, v := range opts {
		for _, iv := range v {
			query.Add(k, iv)
		}
	}
	apiURL.RawQuery = query.Encode()

	req, err := http.NewRequest(http.MethodGet, apiURL.String(), nil)
	if err != nil {
		return localResponse, nil, err
	}

	res, err := c.Do(req)
	if err != nil {
		return localResponse, res, err
	}

	if res.StatusCode != http.StatusOK {
		logger.ErrorKV(ctx, "Bad status code", res.StatusCode)
	}

	data, _ := ioutil.ReadAll(res.Body)
	devices := new(models2.ListItemsResponse)

	err = json.Unmarshal(data, &devices)
	if err != nil {
		return localResponse, res, err
	}
	return *devices, res, nil
}

func (c *client) DescribeDevice(ctx context.Context, deviceID string) (models2.DescribeDeviceResponse, *http.Response, error) {
	var localResponse models2.DescribeDeviceResponse

	apiURLString := c.BasePath + "/api/v1/devices/{deviceId}"
	apiURLString = strings.Replace(apiURLString, "{deviceId}", fmt.Sprintf("%v", deviceID), -1)
	apiURL, err := url.Parse(apiURLString)
	if err != nil {
		return localResponse, nil, err
	}

	req, err := http.NewRequest(http.MethodGet, apiURL.String(), nil)
	if err != nil {
		return localResponse, nil, err
	}

	res, err := c.Do(req)
	if err != nil {
		return localResponse, res, err
	}

	if res.StatusCode != http.StatusOK {
		logger.ErrorKV(ctx, "Bad status code", res.StatusCode)
	}

	data, _ := ioutil.ReadAll(res.Body)
	device := new(models2.DescribeDeviceResponse)

	err = json.Unmarshal(data, &device)
	if err != nil {
		return localResponse, res, err
	}
	return *device, res, nil
}

func (c *client) RemoveDevice(ctx context.Context, deviceID string) (models2.RemovedDevice, *http.Response, error) {
	var localResponse models2.RemovedDevice

	apiURLString := c.BasePath + "/api/v1/devices/{deviceId}"
	apiURLString = strings.Replace(apiURLString, "{deviceId}", fmt.Sprintf("%v", deviceID), -1)
	apiURL, err := url.Parse(apiURLString)
	if err != nil {
		return localResponse, nil, err
	}

	req, err := http.NewRequest(http.MethodDelete, apiURL.String(), nil)
	if err != nil {
		return localResponse, nil, err
	}

	res, err := c.Do(req)
	if err != nil {
		return localResponse, res, err
	}

	if res.StatusCode != http.StatusOK {
		logger.ErrorKV(ctx, "Bad status code", res.StatusCode)
	}

	data, _ := ioutil.ReadAll(res.Body)
	removedDevice := new(models2.RemovedDevice)

	err = json.Unmarshal(data, &removedDevice)
	if err != nil {
		return localResponse, res, err
	}
	return *removedDevice, res, nil
}
