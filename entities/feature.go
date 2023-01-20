package entities

type Feature struct {
	Key             string            `json:"key"`
	Class           *ReferenceItem    `json:"class"`
	Name            string            `json:"name"`
	Level           int               `json:"level"`
	FeatureSpecific *SubFeatureOption `json:"feature_specific"`
	Invocations     []*ReferenceItem  `json:"invocations"`
}

type SubFeatureOption struct {
	SubFeatureOptions *ChoiceOption `json:"subfeature_options"`
}
