package entities

type ProficiencyType string

const (
	ProficiencyTypeArmor       ProficiencyType = "armor"
	ProficiencyTypeWeapon      ProficiencyType = "weapon"
	ProficiencyTypeTool        ProficiencyType = "tool"
	ProficiencyTypeSavingThrow ProficiencyType = "saving-throw"
	ProficiencyTypeSkill       ProficiencyType = "skill"
	ProficiencyTypeInstrument  ProficiencyType = "instrument"
	ProficiencyTypeUnknown     ProficiencyType = ""
)

type Proficiency struct {
	Key       string          `json:"key"`
	Name      string          `json:"name"`
	Type      ProficiencyType `json:"type"`
	Reference *ReferenceItem  `json:"reference"`
}
