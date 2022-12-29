package entities

type Equipment struct {
	Key               string             `json:"key"`
	Name              string             `json:"name"`
	EquipmentCategory *EquipmentCategory `json:"equipment_category"`
	Cost              *Cost              `json:"cost"`
	Weight            int                `json:"weight"`
}

type EquipmentCategory struct {
	Key  string `json:"key"`
	Name string `json:"name"`
}

type Cost struct {
	Quantity int    `json:"quantity"`
	Unit     string `json:"unit"`
}
