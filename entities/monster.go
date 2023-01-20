package entities

type Monster struct {
	Key                   string                `json:"index"`
	Name                  string                `json:"name"`
	Size                  string                `json:"size"`
	Type                  string                `json:"type"`
	Alignment             string                `json:"alignment"`
	ArmorClass            int                   `json:"armor_class"`
	HitPoints             int                   `json:"hit_points"`
	HitDice               string                `json:"hit_dice"`
	Speed                 *Speed                `json:"speed"`
	Strength              int                   `json:"strength"`     //TODO: Refactor to AbilityScore
	Dexterity             int                   `json:"dexterity"`    //TODO: Refactor to AbilityScore
	Constitution          int                   `json:"constitution"` //TODO: Refactor to AbilityScore
	Intelligence          int                   `json:"intelligence"` //TODO: Refactor to AbilityScore
	Wisdom                int                   `json:"wisdom"`       //TODO: Refactor to AbilityScore
	Charisma              int                   `json:"charisma"`     //TODO: Refactor to AbilityScore
	Proficiencies         []*MonsterProficiency `json:"proficiencies"`
	DamageVulnerabilities []string              `json:"damage_vulnerabilities"`
	DamageResistances     []string              `json:"damage_resistances"`
	DamageImmunities      []string              `json:"damage_immunities"`
	ConditionImmunities   []string              `json:"condition_immunities"`
	Senses                *Senses               `json:"senses"`
	Languages             []string              `json:"languages"`
	ChallengeRating       float32               `json:"challenge_rating"`
	XP                    int                   `json:"xp"`
	SpecialAbilities      []*SpecialAbility     `json:"special_abilities"`
	MonsterActions        []*MonsterAction      `json:"actions"`
	LegendaryActions      []*LegendaryAction    `json:"legendary_actions"`
}

type Speed struct {
	Walk   string `json:"walk"`
	Burrow string `json:"burrow"`
	Fly    string `json:"fly"`
	Swim   string `json:"swim"`
	Climb  string `json:"climb"`
}

type MonsterProficiency struct {
	Value       int            `json:"value"`
	Proficiency *ReferenceItem `json:"proficiency"`
}

type Senses struct {
	Blindsight        string `json:"blindsight"`
	Darkvision        string `json:"darkvision"`
	Tremorsense       string `json:"tremorsense"`
	Truesight         string `json:"truesight"`
	PassivePerception int    `json:"passive_perception"`
}

type SpecialAbility struct {
	Name        string   `json:"name"`
	Description []string `json:"desc"`
	Usage       *Usage   `json:"usage"`
}

type Usage struct {
	UsageType      string   `json:"type"`
	UsageTimes     int      `json:"times"`
	UsageRestTypes []string `json:"rest_types"`
}

type MonsterAction struct {
	Name            string     `json:"name"`
	MultiAttackType string     `json:"multiattack_type"`
	AttackBonus     int        `json:"attack_bonus"`
	Description     string     `json:"desc"`
	Actions         []*Actions `json:"actions"`
	Damage          *Damage    `json:"damage"`
}

type Actions struct {
	ActionName  string     `json:"action_name"`
	Count       int        `json:"count"`
	Type        string     `json:"type"`
	AttackBonus int        `json:"attack_bonus"`
	DC          *MonsterDC `json:"dc"`
}

type MonsterDC struct {
	DCType      *ReferenceItem `json:"dc_type"`
	DCValue     int            `json:"dc_value"`
	SuccessType string         `json:"success"`
}

type LegendaryAction struct {
	Name        string     `json:"name"`
	Description string     `json:"desc"`
	DC          *MonsterDC `json:"dc"`
	Damage      *Damage    `json:"damage"`
}
