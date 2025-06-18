package entities

// EquipmentCategory represents a category of equipment (e.g., martial-weapons)
type EquipmentCategory struct {
	Index     string          `json:"index"`
	Name      string          `json:"name"`
	Equipment []*ReferenceItem `json:"equipment"`
	URL       string          `json:"url"`
}