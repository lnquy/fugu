package base

type (
	Struct struct {
		Name   string   `json:"name"`
		Fields []*Field `json:"fields"`
	}

	Field struct {
		Name    string `json:"name"`
		Type    string `json:"type"`
		Size    uint8  `json:"size"`
		Index   uint8  `json:"index"`
		Padding uint8  `json:"padding"`
	}
)
