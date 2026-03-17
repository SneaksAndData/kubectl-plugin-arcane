package publisher

type AllStreamDefinitions struct {
	namespace string
}

func NewAllStreamDefinitionsPublisher(namespace string) *AllStreamDefinitions {
	return &AllStreamDefinitions{
		namespace: namespace,
	}
}
