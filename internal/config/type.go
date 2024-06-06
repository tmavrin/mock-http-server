package config

type Type string

const (
	JSON    Type = "json"
	Swagger Type = "swagger"
)

func (m Type) String() string {
	return string(m)
}

func (m *Type) Parse(s string) {
	switch s {
	case "json":
		*m = JSON
	case "swagger":
		*m = Swagger
	}
}
