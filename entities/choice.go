package entities

type Choice struct {
	Choose    int       `json:"choose"`
	Type      string    `json:"type"`
	OptionSet OptionSet `json:"from"`
}

type ReferenceItem struct {
	Key  string `json:"index"`
	Name string `json:"name"`
	URL  string `json:"url"`
}

type OptionSetType string

const (
	OptionSetTypeArray             OptionSetType = "options_array"
	OptionSetTypeEquipmentCategory OptionSetType = "equipment_category"
)

type OptionSet interface {
	GetType() OptionSetType
}

type OptionsArrayOptionSet struct {
	Options []Option `json:"options"`
}

func (o *OptionsArrayOptionSet) GetType() OptionSetType {
	return OptionSetTypeArray
}

type EquipmentCategoryOptionSet struct {
}

func (o *EquipmentCategoryOptionSet) GetType() OptionSetType {
	return OptionSetTypeEquipmentCategory
}

type OptionType string

const (
	OptionTypeReference        OptionType = "reference"
	OptionTypeChoice           OptionType = "choice"
	OptionTypeCountedReference OptionType = "counted_reference"
)

type Option interface {
	GetType() OptionType
}

type ReferenceOption struct {
	Reference *ReferenceItem `json:"item"`
}

func (o *ReferenceOption) GetType() OptionType {
	return OptionTypeReference
}

type ChoiceOption struct {
	Choice Choice `json:"choice"`
}

func (o *ChoiceOption) GetType() OptionType {
	return OptionTypeChoice
}

type CountedReferenceOption struct {
	CountedReference ReferenceItem `json:"of"`
	Count            int           `json:"count"`
}

func (o *CountedReferenceOption) GetType() OptionType {
	return OptionTypeCountedReference
}
