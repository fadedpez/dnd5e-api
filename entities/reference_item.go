package entities

import "strings"

type ReferenceItem struct {
	Key  string `json:"index"`
	Name string `json:"name"`
	Type string `json:"url"`
}

func (r *ReferenceItem) GetType() string {
	if r.Type == "" {
		return ""
	}

	parts := strings.Split(r.Type, "/")
	if len(parts) < 2 || len(parts) > 4 {
		return ""
	}

	return parts[len(parts)-2]
}
