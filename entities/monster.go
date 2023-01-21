package entities

type Monster struct {
	Key                   string `json:"index"`
	Name                  string `json:"name"`
	Size                  string `json:"size"`
	Type                  string `json:"type"`
	Alignment             string `json:"alignment"`
	MonsterStats          *MonsterStats
	Proficiencies         []*MonsterProficiency `json:"proficiencies"`
	DamageVulnerabilities []string              `json:"damage_vulnerabilities"`
	DamageResistances     []string              `json:"damage_resistances"`
	DamageImmunities      []string              `json:"damage_immunities"`
	ConditionImmunities   []string              `json:"condition_immunities"`
	Senses                *Senses               `json:"senses"`
	Languages             []string              `json:"languages"`
	ChallengeRating       float32               `json:"challenge_rating"`
	XP                    int                   `json:"xp"`
	MonsterActions        []*MonsterAction      `json:"actions"` //TODO: Interface
	MonsterImageURL       string                `json:"image"`
	//TODO: Add legendary actions
	//TODO: Add reactions
	//TODO: Add special abilities
}

type Speed struct {
	Walk   string `json:"walk"`
	Burrow string `json:"burrow"`
	Fly    string `json:"fly"`
	Swim   string `json:"swim"`
	Climb  string `json:"climb"`
}

type MonsterStats struct {
	ArmorClass   int    `json:"armor_class"`
	HitPoints    int    `json:"hit_points"`
	HitDice      string `json:"hit_dice"`
	Speed        *Speed `json:"speed"`
	Strength     int    `json:"strength"` //TODO: Refactor to AbilityScore
	Dexterity    int    `json:"dexterity"`
	Constitution int    `json:"constitution"`
	Intelligence int    `json:"intelligence"`
	Wisdom       int    `json:"wisdom"`
	Charisma     int    `json:"charisma"`
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
	Name        string  `json:"name"`
	AttackBonus int     `json:"attack_bonus"`
	Description string  `json:"desc"`
	Damage      *Damage `json:"damage"`
}
