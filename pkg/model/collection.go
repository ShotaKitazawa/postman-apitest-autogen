package model

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// CollectionBase

type CollectionBase struct {
	Info CollectionInfo   `json:"info"`
	Item []CollectionItem `json:"item"` // TODO: support ItemGroup (aka. folder)
	// untouch fields
	Event                   json.RawMessage `json:"event"`
	Variable                json.RawMessage `json:"variable"`
	Auth                    json.RawMessage `json:"auth"`
	ProtocolProfileBehavior json.RawMessage `json:"protocolProfileBehavior"`
}

func NewCollectionFromFile(filepath string) (*CollectionBase, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	var c CollectionBase
	if err := json.Unmarshal(data, &c); err != nil {
		return nil, err
	}
	if !c.IsValidVersion() {
		return nil, fmt.Errorf("schema is unknown version")
	}
	return &c, nil
}

func (b CollectionBase) IsValidVersion() bool {
	var expectedVersions = []string{
		"https://schema.getpostman.com/json/collection/v2.1.0/collection.json",
	}
	// contain?
	for _, e := range expectedVersions {
		if b.Info.Schema == e {
			return true
		}
	}
	return false
}

// CollectionInfo

type CollectionInfo struct {
	Name   string `json:"name"`
	Schema string `json:"schema"`
	// untouch fields
	PostmanID   json.RawMessage `json:"_postman_id"`
	Description json.RawMessage `json:"description"`
	Version     json.RawMessage `json:"version"`
}

// CollectionItem

type CollectionItem struct {
	Name      string               `json:"name"`
	Events    []CollectionEvent    `json:"event,omitempty"`
	Responses []CollectionResponse `json:"response"`
	// untouch fields
	ID                      json.RawMessage `json:"id"`
	Description             json.RawMessage `json:"description"`
	Variable                json.RawMessage `json:"variable"`
	Request                 json.RawMessage `json:"request,omitempty"`
	ProtocolProfileBehavior json.RawMessage `json:"protocolProfileBehavior"`
}

func (i CollectionItem) IsJsonResponse() bool {
	for _, r := range i.Responses {
		for _, h := range r.Headers {
			if h.Key == "content-type" && strings.Contains(h.Value.(string), "application/json") {
				return true
			}
		}
	}
	return false
}

func (i CollectionItem) IsOnlyOnceResponse() bool {
	return len(i.Responses) == 1
}

// CollectionEvent

type CollectionEvent struct {
	Listen string `json:"listen"`
	Script struct {
		Execes []string `json:"exec"`
		Type   string   `json:"type"`
	} `json:"script"`
}

func NewCollectionEvent(test string) CollectionEvent {
	return CollectionEvent{
		Listen: "test",
		Script: struct {
			Execes []string `json:"exec"`
			Type   string   `json:"type"`
		}{
			Execes: []string{test},
			Type:   "text/javascript",
		},
	}
}

// CollectionResponse

type CollectionResponse struct {
	Name    string         `json:"name"`
	Headers []CollectionKV `json:"header"`
	Body    string         `json:"body"`
	Code    int            `json:"code"`
	// untouch fields
	ID                     json.RawMessage `json:"id"`
	OriginalRequest        json.RawMessage `json:"originalRequest"`
	ResponseTime           json.RawMessage `json:"responseTime"`
	Timings                json.RawMessage `json:"timings"`
	Cookies                json.RawMessage `json:"cookie"`
	Status                 string          `json:"status"`
	PostmanPreviewlanguage string          `json:"_postman_previewlanguage"`
}

func (c CollectionResponse) GetBody() []byte {
	return json.RawMessage([]byte(c.Body))
}

type CollectionKV struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}
