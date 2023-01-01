package entities

type Spell struct {
	Key           string           `json:"key"`
	Name          string           `json:"name"`
	Range         string           `json:"range"`
	Ritual        bool             `json:"ritual"`
	Duration      string           `json:"duration"`
	Concentration bool             `json:"concentration"`
	CastingTime   string           `json:"casting_time"`
	SpellLevel    int              `json:"level"`
	SpellDamage   *SpellDamage     `json:"damage"`
	DC            *DC              `json:"dc"`
	AreaOfEffect  *AreaOfEffect    `json:"area_of_effect"`
	SpellSchool   *ReferenceItem   `json:"school"`
	SpellClasses  []*ReferenceItem `json:"classes"`
}

type SpellDamage struct {
	SpellDamageType        *ReferenceItem          `json:"damage_type"`
	SpellDamageAtSlotLevel *SpellDamageAtSlotLevel `json:"damage_at_slot_level"`
}

type SpellDamageAtSlotLevel struct {
	FirstLevel   string `json:"1"`
	SecondLevel  string `json:"2"`
	ThirdLevel   string `json:"3"`
	FourthLevel  string `json:"4"`
	FifthLevel   string `json:"5"`
	SixthLevel   string `json:"6"`
	SeventhLevel string `json:"7"`
	EighthLevel  string `json:"8"`
	NinthLevel   string `json:"9"`
}

type DC struct {
	DCType    *DCType `json:"dc_type"`
	DCSuccess string  `json:"dc_success"`
}

type DCType struct {
	Key  string `json:"key"`
	Name string `json:"name"`
}

type AreaOfEffect struct {
	Type string `json:"type"`
	Size int    `json:"size"`
}
