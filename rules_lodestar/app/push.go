package main

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	os    = runtime.GOOS
	sm    = make(map[string]string)
	paths = []string{}
)

func main() {

	filepath.WalkDir(".", getKeyValues)

	if len(paths) == 0 {
		log.Fatalf("could not get paths")
	}

	for _, i := range paths {
		content, err := ioutil.ReadFile(i)
		if err != nil {
			log.Fatalf("could not read file %s: %s", i, err)
		}

		lines := strings.Split(string(content), "\n")

		for _, line := range lines {
			kv := strings.Split(line, " ")

			if len(kv) <= 1 {
				continue
			}

			sm[kv[0]] = kv[1]
		}
	}

	var lp string
	var acp string

	switch os {
	case "windows":
		lp = filepath.Join("\\", sm["STABLE_WORKSPACE_DIR"], sm["STABLE_LODESTAR_DIR"])
		acp = filepath.Join("\\", sm["STABLE_WORKSPACE_DIR"], sm["STABLE_APPCONFIG_DIR"])
	default:
		lp = filepath.Join("/", sm["STABLE_WORKSPACE_DIR"], sm["STABLE_LODESTAR_DIR"])
		acp = filepath.Join("/", sm["STABLE_WORKSPACE_DIR"], sm["STABLE_APPCONFIG_DIR"])
	}

	e := sm["STABLE_ENVIRONMENT"]
	yk := sm["STABLE_YAML_KEYS"]
	t := sm["STABLE_GIT_TOKEN"]
	u := sm["STABLE_GIT_USER"]

	out, err := exec.Command(lp, "app", "push", "--env", e, "--config-path", acp, "--yaml-keys", yk, "--username", u, "--token", t).Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", out)

}

func getKeyValues(s string, d fs.DirEntry, e error) error {
	if e != nil {
		return e
	}
	if !d.IsDir() {
		switch os {
		case "windows":
			//to do: get this to work. Paths are not being picked up during crawl.  May be a bazel thing but I dont have
			//a windows machine to test appropriately
			if strings.Contains(s, "stamps.txt") || strings.Contains(s, "lodestar.txt") {
				paths = append(paths, s)
				fmt.Println(s)
			}
		default:
			if strings.Contains(s, "/stamps.txt") || strings.Contains(s, "/lodestar.txt") {
				paths = append(paths, s)
			}
		}
	}
	return nil
}
