package utils

const (
	prefix = "github.com/serpro69/pgc"
)

// ResourceType is a custom type to represent the type of a project component resource
type ResourceType struct {
	p string
	t string
}

func (r ResourceType) String() string {
	return prefix + ":" + r.p + ":" + r.t
}

// NewComponentType returns a new instance of ComponentResourceType
func NewResourceType(p, t string) ResourceType {
	return ResourceType{
		p: p,
		t: t,
	}
}
