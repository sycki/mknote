package structs

type Model struct {
	Map map[string]interface{}
}

func (m *Model) Set(key string, v interface{}) {
	m.Map[key] = v
}

func (m *Model) Get(key string) interface{} {
	return m.Map[key]
}

func (m *Model) Clear() {
	for k := range m.Map {
		delete(m.Map, k)
	}
}
