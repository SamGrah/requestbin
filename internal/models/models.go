package models

import "time"

type Request struct {
	Id         int64
	RecievedAt time.Time
	Headers    string
	Body       string
	Host       string
	Method     string
	Bin        string
}

type Bin struct {
	BinId     string
	CreatedAt time.Time
	Owner     string
}
