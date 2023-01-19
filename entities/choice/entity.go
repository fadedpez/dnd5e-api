package choice

type OptionType string

const (
	OptionTypeReference          OptionType = "reference"
	OptionTypeChoice             OptionType = "choice"
	OptionalTypeCountedReference OptionType = "counted_reference"
	OptionTypeMultiple           OptionType = "multiple"
)

type Option interface {
	GetOptionType() OptionType
}

type ReferenceItem struct {
	Key  string `json:"index"`
	Name string `json:"name"`
	Type string `json:"type"` //TODO: make this an enum
}

type CountedReferenceOption struct {
	Count     int            `json:"count"`
	Reference *ReferenceItem `json:"reference"`
}

func (o *CountedReferenceOption) GetOptionType() OptionType {
	return OptionalTypeCountedReference
}

type ReferenceOption struct {
	Reference *ReferenceItem `json:"reference"`
}

func (o *ReferenceOption) GetOptionType() OptionType {
	return OptionTypeReference
}

type ChoiceOption struct {
	Description string      `json:"description"`
	ChoiceCount int         `json:"choice_count"`
	ChoiceType  string      `json:"choice_type"`
	OptionList  *OptionList `json:"option_list"`
}

func (o *ChoiceOption) GetOptionType() OptionType {
	return OptionTypeChoice
}

type MultipleOption struct {
	Items []Option `json:"items"`
}

func (o *MultipleOption) GetOptionType() OptionType {
	return OptionTypeMultiple
}

type OptionList struct {
	Options []Option `json:"option_list"`
}
