package entities

type Choice struct {
	Choose    int       `json:"choose"`
	Type      string    `json:"type"`
	OptionSet *OptionList `json:"from"`
	chosen int // number of options chosen
}

type OptionSetType string

const (
	OptionSetTypeArray             OptionSetType = "options_array"
)

type OptionList struct {
	Options []Option `json:"options"`
}

func (o *OptionList) GetType() OptionSetType {
	return OptionSetTypeArray
}

type OptionType string

const (
	OptionTypeReference        OptionType = "reference"
	OptionTypeChoice           OptionType = "choice"
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

func (o *Choice) GetType() OptionType {
	return OptionTypeChoice
}


