package models

// CreateDeviceResponse is the response body from the CreateDevice endpoint
type CreateDeviceResponse struct {
	DeviceID int `json:"deviceId,string"`
}

// DescribeDeviceResponse is the response of DescribeDeviceRequest
type DescribeDeviceResponse struct {
	Value Item `json:"value"`
}

// CreateDeviceRequest is the request struct for api CreateDevice
type CreateDeviceRequest struct {
	Platform string `json:"platform"`
	UserID   string `json:"userId"`
}

// RemovedDevice Response is the response of DeletedDeviceRequest
type RemovedDevice struct {
	Found bool `json:"found"`
}
