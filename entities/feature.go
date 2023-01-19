package entities

type Feature struct {
	Key             string         `json:"key"`
	Class           *ReferenceItem `json:"class"`
	Name            string         `json:"name"`
	Level           int            `json:"level"`
	FeatureSpecific *SubFeature    `json:"feature_specific"`
}

type SubFeature struct {
	SubfeatureOptions *ChoiceOption    `json:"subfeature_options"`
	Invocations       *[]ReferenceItem `json:"invocations"`
}
