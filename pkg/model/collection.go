package model

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

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

type CollectionBase struct {
	Info CollectionInfo   `json:"info"`
	Item []CollectionItem `json:"item"`
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

type CollectionInfo struct {
	PostmanID string `json:"_postman_id"`
	Name      string `json:"name"`
	Schema    string `json:"schema"`
}

type CollectionItem struct {
	Name      string               `json:"name"`
	Events    []CollectionEvent    `json:"event,omitempty"`
	Request   CollectionRequest    `json:"request,omitempty"`
	Responses []CollectionResponse `json:"response"`
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

type CollectionEvent struct {
	Listen string           `json:"listen"`
	Script CollectionScript `json:"script"`
}

func NewCollectionEvent(test string) CollectionEvent {
	return CollectionEvent{
		Listen: "test",
		Script: CollectionScript{
			Execes: []string{test},
			Type:   "text/javascript",
		},
	}
}

type CollectionScript struct {
	Execes []string `json:"exec"`
	Type   string   `json:"type"`
}

type CollectionRequest struct {
	Method  string                `json:"method"`
	Headers []CollectionKV        `json:"header"`
	URL     CollectionURL         `json:"url"`
	Body    CollectionRequestBody `json:"body,omitempty"`
}

type CollectionRequestBody struct {
	Mode string `json:"mode"`
	Raw  string `json:"raw"`
}

type CollectionResponse struct {
	Name                   string             `json:"name"`
	OriginalRequest        CollectionRequest  `json:"originalRequest"`
	Status                 string             `json:"status"`
	Code                   int                `json:"code"`
	PostmanPreviewlanguage string             `json:"_postman_previewlanguage"`
	Headers                []CollectionKV     `json:"header"`
	Cookies                []CollectionCookie `json:"cookie"`
	Body                   string             `json:"body"`
}

func (c CollectionResponse) GetBody() []byte {
	return json.RawMessage([]byte(c.Body))
}

type CollectionURL struct {
	Raw      string         `json:"raw"`
	Protocol string         `json:"protocol"`
	Host     []string       `json:"host"`
	Port     string         `json:"port"`
	Path     []string       `json:"path"`
	Query    []CollectionKV `json:"query,omitempty"`
}

type CollectionURLQuery struct {
}

type CollectionCookie map[string]string

type CollectionKV struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}
