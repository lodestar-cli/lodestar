package remote

import (
	"bytes"
	"testing"
)

func TestManagementFile_GetByteContent(t *testing.T) {
	m := getTestManagementFile()

	b := m.GetByteContent()

	if !bytes.Equal(m.ByteContent, b) {
		t.Error("Error Updating File: byte content does not update")
	}
}

func TestManagementFile_GetStringContent(t *testing.T) {
	m := getTestManagementFile()

	s := m.GetStringContent()

	if m.StringContent != s {
		t.Error("Error Updating File: string content does not update")
	}
}

func TestManagementFile_GetKeyValues(t *testing.T) {
	m := getTestManagementFile()

	testTable := []struct {
		KeyList     []string
		ExpectedMap map[string]string
		ExpectError bool
	}{
		{[]string{"test", "test1"}, map[string]string{"test": "fail", "test1": "fail1"}, false},
		{[]string{"test", "test5"}, nil, true},
	}

	for _, test := range testTable {
		kv, err := m.GetKeyValues(test.KeyList)
		if err != nil {
			if test.ExpectError {
				continue
			}
			t.Errorf("Error Getting Keys: %s", err)
		}

		if test.ExpectError && err == nil {
			t.Error("Expected GetKeyValues to fail but it passed")
		}

		if !test.ExpectError {
			for k, _ := range kv {
				if kv[k] != test.ExpectedMap[k] {
					t.Errorf("Error after updating EnvironmentStateGraph: for key %s: expected %s but got %s", k, test.ExpectedMap[k], kv[k])
				}
			}
		}

	}
}

func TestManagementFile_UpdateFileContents(t *testing.T) {
	m := getTestManagementFile()

	expectedPass := `
test: "pass"
  test1: "pass1"
test2: fail2
`

	testTable := []struct {
		KeyValueMap   map[string]string
		ExpectString  string
		ExpectByte    []byte
		ExpectUpdated bool
		ExpectError   bool
	}{
		{map[string]string{"test": "pass", "test1": "pass1"}, expectedPass, []byte(expectedPass), true, false},
		{map[string]string{"test": "fail", "test1": "fail1"}, m.StringContent, m.ByteContent, false, false},
	}

	// test update
	for _, test := range testTable {
		m = getTestManagementFile()
		u, err := m.UpdateFileContents(test.KeyValueMap)
		if err != nil {
			if test.ExpectError {
				continue
			}
			t.Errorf("Error Updating File Content: %s", err)
		}

		if test.ExpectError && err == nil {
			t.Error("Expected UpdateFileContents to fail but it passed")
		}

		if u != test.ExpectUpdated {
			t.Errorf("Bad Update: Expected %v and got %v", test.ExpectUpdated, u)
		}

		if !test.ExpectError {
			if m.StringContent != test.ExpectString {
				t.Error("Error Updating File: string content does not update")
			}
			if !bytes.Equal(m.ByteContent, test.ExpectByte) {
				t.Error("Error Updating File: byte content does not update")
			}
		}

	}

	// test for bad string content
	m = getTestManagementFile()
	m.StringContent = `
test: fail
test1: ""
test2: fail2
`
	_, err := m.UpdateFileContents(map[string]string{"test": "", "test1": "pass1"})
	if err == nil {
		t.Error("Expected UpdateFileContents to fail due to empty string value but it passed")
	}

}

func getTestManagementFile() *ManagementFile {
	testYaml := `
test: "fail"
  test1: fail1
test2: fail2
`
	m := ManagementFile{
		Name:          "test",
		Path:          "/test/test.yaml",
		ByteContent:   []byte(testYaml),
		StringContent: testYaml,
	}

	return &m
}
