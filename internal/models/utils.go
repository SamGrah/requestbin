package models

import (
	"bytes"
	"encoding/gob"
	"time"
)

func TimeToString(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

func TimeFromString(t string) (time.Time, error) {
	return time.Parse("2006-01-02 15:04:05", t)
}

func encodeMapToString(m map[string][]string) (string, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(m); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func decodeStringToMap(s string) (map[string][]string, error) {
	var m map[string][]string
	dec := gob.NewDecoder(bytes.NewBufferString(s))
	if err := dec.Decode(&m); err != nil {
		return nil, err
	}
	return m, nil
}
