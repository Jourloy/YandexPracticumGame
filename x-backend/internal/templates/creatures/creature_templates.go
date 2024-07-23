package creature_templates

type template struct {
	MaxStorage         int     `json:"maxStorage"`
	UsedStorage        int     `json:"usedStorage"`
	MaxHealth          int     `json:"maxHealth"`
	Health             int     `json:"health"`
	FatiguePerStep     float64 `json:"fatiguePerStep"`
	FatigueModificator float64 `json:"fatigueModificator"`
	Fatigue            float64 `json:"fatigue"`
	RequireFood        float64 `json:"requireFood"`
}

var CreatureTemplate = map[string]template{}

func init() {
	CreatureTemplate[`human`] = template{
		MaxStorage:         50,
		UsedStorage:        0,
		FatiguePerStep:     0,
		FatigueModificator: 1,
		Fatigue:            0,
		MaxHealth:          100,
		Health:             100,
		RequireFood:        1,
	}

	CreatureTemplate[`swarm`] = template{
		MaxStorage:         50,
		UsedStorage:        0,
		FatiguePerStep:     0,
		FatigueModificator: 0.5,
		Fatigue:            0,
		MaxHealth:          50,
		Health:             50,
		RequireFood:        0.5,
	}

	CreatureTemplate[`robot`] = template{
		MaxStorage:         50,
		UsedStorage:        0,
		FatiguePerStep:     2,
		FatigueModificator: 2,
		Fatigue:            0,
		MaxHealth:          200,
		Health:             100,
		RequireFood:        1,
	}
}
