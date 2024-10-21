package utils

// ResourceType is a custom type to represent the type of a project component resource
type ResourceType string

const (
	prefix   = "github.com/serpro69/pgc:project:"
	Project  = ResourceType(prefix + "Project")
	Services = ResourceType(prefix + "Services")
)

func (r ResourceType) String() string {
	return string(r)
}
