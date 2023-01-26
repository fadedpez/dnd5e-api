package dnd5e

import (
	"strings"

	"github.com/fadedpez/dnd5e-api/entities"
)

type optionSet struct {
	OptionSetType     string         `json:"option_set_type"`
	EquipmentCategory *referenceItem `json:"equipment_category"`
	Options           []*option      `json:"options"`
}

func (o *optionSet) toEntity() *entities.OptionList {
	var options []entities.Option

	for _, opt := range o.Options {
		options = append(options, opt.toEntity())
	}

	return &entities.OptionList{
		Options: options,
	}
}

type option struct {
	OptionType string         `json:"option_type"`
	Count      int            `json:"count"`
	Of         *referenceItem `json:"of"`
	Items      []*option      `json:"items"`
	Item       *referenceItem `json:"item"`
	Choice     *choiceResult  `json:"choice"`
}

func (o *option) toEntity() entities.Option {
	switch o.OptionType {
	case "reference":
		return &entities.ReferenceOption{
			Reference: &entities.ReferenceItem{
				Key:  o.Item.Index,
				Name: o.Item.Name,
				Type: urlToType(o.Item.URL),
			},
		}
	case "choice":
		return o.Choice.toEntity()
	case "counted_reference":
		return &entities.CountedReferenceOption{
			Count: o.Count,
			Reference: &entities.ReferenceItem{
				Key:  o.Of.Index,
				Name: o.Of.Name,
				Type: urlToType(o.Of.URL),
			},
		}
	case "multiple":
		var items []entities.Option
		for _, item := range o.Items {
			items = append(items, item.toEntity())
		}

		return &entities.MultipleOption{
			Items: items,
		}
	}

	return nil
}

type referenceOption struct {
	Item *referenceItem `json:"item"`
}

func urlToType(url string) string {
	if url == "" {
		return ""
	}

	urlparts := strings.Split(url, "/")
	if len(urlparts) < 3 {
		return ""
	}

	return urlparts[2]
}

type countedReferenceOption struct {
	Count int            `json:"count"`
	Of    *referenceItem `json:"of"`
}

type choiceResult struct {
	Desc   string     `json:"desc"`
	Choose int        `json:"choose"`
	Type   string     `json:"type"`
	From   *optionSet `json:"from"`
}

func (c *choiceResult) toEntity() *entities.ChoiceOption {

	return &entities.ChoiceOption{
		Description: c.Desc,
		ChoiceCount: c.Choose,
		ChoiceType:  c.Type,
		OptionList:  c.From.toEntity(),
	}
}

type equipmentCategoryOptionSet struct {
	EquipmentCategory *referenceItem `json:"equipment_category"`
}
