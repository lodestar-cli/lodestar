package environment

type EnvironmentState struct {
	Name string `yaml:"name"`
	YamlKeys  map[string]string `yamlKeys:"tags,omitempty"`
}

//Updates current State tags to new given tags.  If a new tag isn't in the current State graph, it adds it
func (e *EnvironmentState) UpdateKeys(keys map[string]string) bool {
	update := false

	for k,v := range keys{
		val, ok := e.YamlKeys[k]
		if ok {
			if v != val{
				e.YamlKeys[k] = v
				update = true
			}
		} else{
			e.YamlKeys[k] = v
		}
	}

	return update
}
