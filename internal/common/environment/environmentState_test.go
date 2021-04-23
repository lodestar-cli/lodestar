package environment

import "testing"

func TestEnvironmentState_UpdateKeys(t *testing.T) {
	m := map[string]string{
		"tag": "wow!",
	}

	testTable := []struct {
		KeyMap        map[string]string
		ExpectUpdated bool
	}{
		{map[string]string{"tag": "test"}, false},
		{m, true},
	}

	for _, test := range testTable {

		e := EnvironmentState{
			Name:     "dev",
			YamlKeys: map[string]string{"tag": "test"},
		}

		u := e.UpdateKeys(test.KeyMap)

		if u != test.ExpectUpdated {
			t.Errorf("Bad Update: Expected %v and got %v", test.ExpectUpdated, u)
		}
	}
}
