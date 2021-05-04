package home

import (
	"github.com/lodestar-cli/lodestar/internal/cli/home/mocks"
	"os/user"
	"path"
	"testing"
)

func TestGetPath(t *testing.T) {
	usr, err := user.Current()
	if err != nil {
		t.Fatalf("could not get current user: %s", err)
	}

	dir := usr.HomeDir

	testCommand := "test"
	testCommandName := "test"

	lodestarDirectory := path.Join(dir, ".lodestar/", testCommand, testCommandName+".yaml")

	testPath, err := GetPath(testCommand, testCommandName)
	if err != nil {
		t.Fatalf("failed to get path: %s", err)
	}

	if lodestarDirectory != testPath {
		t.Fatalf("paths do not match.  Expected %s but got %s", lodestarDirectory, testPath)
	}
}

func TestGetContent(t *testing.T) {
	iou := new(mocks.IoUtil)

	usr, err := user.Current()
	if err != nil {
		t.Fatalf("could not get current user: %s", err)
	}

	dir := usr.HomeDir

	testCommand := "test"
	testCommandName := "test"

	p := path.Join(dir, ".lodestar/", testCommand, testCommandName+".yaml")
	expectedContent := "content"

	iou.On("ReadFile", p).Return([]byte(expectedContent), nil)
	content, err := GetContent(p, iou)
	if err != nil {
		t.Fatalf("failed getting file content: %s", err)
	}

	if expectedContent != string(content) {
		t.Fatalf("content returned doesn't match. Expected %s but got %s", expectedContent, content)
	}

	iou.AssertExpectations(t)

}
