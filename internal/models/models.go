package models

import (
	"time"
)

type Bin struct {
	BinId     int64
	CreatedAt time.Time
	Owner     string
}

type Request struct {
	Id         int64
	RecievedAt time.Time
	Headers    string
	RemoteAddr string
	Body       string
	Host       string
	RequestUri string
	Method     string
	Bin        int64
}

func (r *Request) GetHeaders() (map[string][]string, error) {
	return decodeStringToMap(r.Headers)
}

func (r *Request) SetHeaders(reqHeaders map[string][]string) error {
	headers, err := encodeMapToString(reqHeaders)
	if err != nil {
		return err
	}
	r.Headers = headers
	return nil
}
