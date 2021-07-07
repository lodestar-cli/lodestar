// +build mage

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strings"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

var (
	env = map[string]string{
		"CGO_ENABLED": "0",
	}
)

func gitCommit() string {
	s, e := sh.Output("git", "rev-parse", "--short", "HEAD")
	if e != nil {
		fmt.Printf("Failed to get GIT version: %s\n", e)
		return ""
	}
	return s
}

func previousMerge() string {
	p, e := sh.Output("git", "log", "--format=format:%H", "--merges", "--skip", "1", "-n", "1")
	if e != nil {
		fmt.Printf("Failed to get previous merge: %s\n", e)
		return ""
	}
	return p
}

// returns a list of files that changed between merges
func ListDependencyChanges() error {
	bazeliskFiles, err := getDependencyChanges()
	if err != nil {
		return fmt.Errorf("could not get dependency changes: %s", err)
	}

	fmt.Println("Files to consider:")
	for _, file := range bazeliskFiles {
		fmt.Println(file)
	}

	return nil
}

// Gets the code coverage as a text file
func Coverage() error {
	err := sh.Run("bazelisk", "coverage", "--combined_report=lcov", "//...")
	if err != nil {
		return fmt.Errorf("failed to run go tests: %s", err)
	}

	l, err := sh.Output("cat", "bazel-out/_coverage/lcov_files.tmp")
	if err != nil {
		return fmt.Errorf("failed to get coverage files: %s", err)
	}

	lcovs := strings.Split(l, "\n")

	var cf string
	for _, lcov := range lcovs {
		cov, err := sh.Output("cat", lcov)
		if err != nil {
			return fmt.Errorf("failed to read coverage file %s: %s", lcov, err)
		}

		s := strings.Split(cov, "\n")
		a := s[1:]

		c := strings.Join(a, "\n")

		cf += "\n" + c
	}

	bcf := []byte(cf)
	bcf = bytes.Trim(bcf, " \n")
	err = ioutil.WriteFile("coverage.txt", bcf, 0644)
	if err != nil {
		return fmt.Errorf("could not write coverage file: %s", err)
	}

	return nil
}

func getDependencyChanges() ([]string, error) {
	commitRange := fmt.Sprintf("%s..%s", previousMerge(), gitCommit())

	fmt.Printf("commit range: %s\n", commitRange)

	changedFiles, e := sh.Output("git", "diff", "--name-only", commitRange)
	if e != nil {
		return nil, fmt.Errorf("Failed to get changeset from git: %w", e)
	}

	bazeliskFiles := []string{}
	for _, file := range strings.Split(changedFiles, "\n") {

		fileType := file[len(file)-3:]

		if fileType != ".go" || strings.Contains(file, "magefile.go") {
			continue
		}

		files, e := sh.Output("bazelisk", "query", file)
		if e != nil {
			fmt.Printf("Unable to query bazelisk with file: %q, %v", file, e)
		} else {
			bazeliskFiles = append(bazeliskFiles, strings.Split(files, "\n")...)
		}
	}

	return bazeliskFiles, nil
}

// runs bazel rules made for lodestar
func RunRules() error {
	o := runtime.GOOS
	commit := gitCommit()
	commit = commit[0:7]
	yk := "tag=" + commit
	env["GIT_USER"] = os.Getenv("GIT_USER")
	env["GIT_TOKEN"] = os.Getenv("GIT_TOKEN")
	env["YAML_KEYS"] = yk

	switch o {
	case "windows":
		err := sh.Run("bazelisk", "build", "//rules_lodestar/app:push")
		if err != nil {
			return fmt.Errorf("failed to run go tests: %s", err)
		}
	default:
		err := sh.RunWith(env, "bazelisk", "run", "//rules_lodestar/app:push")
		if err != nil {
			return fmt.Errorf("failed to run go tests: %s", err)
		}
	}

	return nil
}

type CI mg.Namespace

// updates Docker images based upon the changes to the code base
func (CI) Build() error {
	bazeliskFiles, err := getDependencyChanges()
	if err != nil {
		return fmt.Errorf("could not get dependency changes: %s", err)
	}

	fmt.Println("Files to consider:")
	for _, file := range bazeliskFiles {
		fmt.Println(file)
	}

	fileSplat := strings.Join(bazeliskFiles, " ")

	binaries, e := sh.Output("bazelisk", "query", "--keep_going", "--noshow_progress",
		fmt.Sprintf("kind(.*_binary, rdeps(//..., set(%s)))", fileSplat))
	if e != nil {
		return fmt.Errorf("Failed to get list of changed binaries: %w", e)
	}

	images, e := sh.Output("bazelisk", "query", "--keep_going", "--noshow_progress",
		fmt.Sprintf("kind(container_image, rdeps(//..., set(%s)))", fileSplat))
	if e != nil {
		return fmt.Errorf("Failed to get list of changed images: %w", e)
	}

	dockerPushAmd, e := sh.Output("bazelisk", "query", "--keep_going", "--noshow_progress",
		fmt.Sprintf("filter(amd, kind(container_push, rdeps(//..., set(%s))))", fileSplat))
	if e != nil {
		return fmt.Errorf("Failed to list of containers that need to be pushed: %w", e)
	}

	dockerPushArm, e := sh.Output("bazelisk", "query", "--keep_going", "--noshow_progress",
		fmt.Sprintf("filter(arm, kind(container_push, rdeps(//..., set(%s))))", fileSplat))
	if e != nil {
		return fmt.Errorf("Failed to list of containers that need to be pushed: %w", e)
	}

	fmt.Printf("Files to consider: \nBinaries: %s\nImages: %s\n, Docker: %s\n", binaries, images, dockerPushAmd)

	args := []string{}

	err = bazeliskRun("binary", "build", binaries, args...)
	if err != nil {
		return err
	}

	err = bazeliskRun("images", "build", images, args...)
	if err != nil {
		return err
	}

	args = append(args, "--platforms=@io_bazel_rules_go//go/toolchain:linux_amd64")

	err = bazeliskRun("pushImage", "run", dockerPushAmd, args...)
	if err != nil {
		return fmt.Errorf("failed pushing arm image to docker: %s", err)
	}

	args = []string{}
	args = append(args, "--platforms=@io_bazel_rules_go//go/toolchain:linux_arm64")

	err = bazeliskRun("pushImage", "run", dockerPushArm, args...)
	if err != nil {
		return fmt.Errorf("failed pushing arm image to docker: %s", err)
	}

	return nil
}

func bazeliskRun(name string, cmd string, files string, args ...string) error {
	if files == "" {
		fmt.Printf("No %s to build\n", name)
	} else {
		fmt.Printf("Building %s...\n", name)

		for _, v := range strings.Split(files, "\n") {

			a := []string{
				cmd,
			}
			for _, arg := range args {
				a = append(a, arg)
			}
			a = append(a, v)
			err := sh.Run("bazelisk", a...)
			if err != nil {
				return fmt.Errorf("Failed to run bazelisk cmd: %s on target: %s -%w", cmd, v, err)
			}
		}
	}
	return nil
}
