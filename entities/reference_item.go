package entities

import "strings"

type ReferenceItem struct {
	Key  string `json:"index"`
	Name string `json:"name"`
	URL  string `json:"url"`
}

func (r *ReferenceItem) GetType() string {
	if r.URL == "" {
		return ""
	}

	parts := strings.Split(r.URL, "/")
	if len(parts) < 2 || len(parts) > 4 {
		return ""
	}

	return parts[len(parts)-2]
}