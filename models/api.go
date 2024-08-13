package models

type ApiResp[T any] struct {
	Payload T      `json:"payload,omitempty"`
	Err     string `json:"error,omitempty"`
}
