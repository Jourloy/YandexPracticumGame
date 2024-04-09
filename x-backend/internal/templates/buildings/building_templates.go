package building_templates

type template struct {
	MaxDurability int  `json:"maxDurability"`
	Durability    int  `json:"durability"`
	MaxStorage    int  `json:"maxStorage"`
	UsedStorage   int  `json:"usedStorage"`
	Level         int  `json:"level"`
	AttackRange   int  `json:"attackRange"`
	CanTrade      bool `json:"canTrade"`
}

var BuildingTemplates = map[string]template{}

func init() {
	BuildingTemplates[`townhall`] = template{
		MaxDurability: 1000,
		Durability:    1000,
		MaxStorage:    200,
		UsedStorage:   0,
		Level:         0,
		AttackRange:   10,
		CanTrade:      true,
	}

	BuildingTemplates[`tower`] = template{
		MaxDurability: 500,
		Durability:    500,
		MaxStorage:    100,
		UsedStorage:   0,
		Level:         0,
		AttackRange:   20,
		CanTrade:      false,
	}

	BuildingTemplates[`wall`] = template{
		MaxDurability: 500,
		Durability:    500,
		MaxStorage:    0,
		UsedStorage:   0,
		Level:         0,
		AttackRange:   0,
		CanTrade:      false,
	}

	BuildingTemplates[`storage`] = template{
		MaxDurability: 500,
		Durability:    500,
		MaxStorage:    1000,
		UsedStorage:   0,
		Level:         0,
		AttackRange:   0,
		CanTrade:      false,
	}

	BuildingTemplates[`market`] = template{
		MaxDurability: 500,
		Durability:    500,
		MaxStorage:    1000,
		UsedStorage:   0,
		Level:         0,
		AttackRange:   0,
		CanTrade:      true,
	}
}
