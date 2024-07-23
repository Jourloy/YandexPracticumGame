package item_templates

type template struct {
	Weight int `json:"weight"`
}

var ItemTemplates = map[string]template{}

func init() {
	ItemTemplates[`topor`] = template{
		Weight: 5,
	}
}
