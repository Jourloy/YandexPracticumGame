package resource_templates

type template struct {
	Weight int `json:"weight"`
}

var ResourceTemplates = map[string]template{}

func init() {
	ResourceTemplates[`log`] = template{
		Weight: 10,
	}
}
