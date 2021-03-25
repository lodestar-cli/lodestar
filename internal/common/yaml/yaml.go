package yaml

import (
	"errors"
	"strings"
)

func ReplaceTag(yaml string, tag string) (string, error) {
	split := strings.Split(yaml, "\n")

	tagLines, err := findTag(split)
	if err != nil {
		return "", err
	}

	l := tagLines[0] + 1
	tagSplit := strings.Split(split[l], " ")

	for i, txt := range tagSplit {
		if strings.Contains(txt, "tag:") {
			tagSplit[i+1] = tag
			break
		}
	}
	split[l] = strings.Join(tagSplit, " ")
	yaml = strings.Join(split, "\n")

	return yaml, nil
}

func GetTag(yaml string) (string, error) {
	var tag string
	split := strings.Split(yaml, "\n")

	tagLines, err := findTag(split)
	if err != nil {
		return "", err
	}

	l := tagLines[0] + 1
	tagSplit := strings.Split(split[l], " ")

	for i, txt := range tagSplit {
		if strings.Contains(txt, "tag:") {
			tag = tagSplit[i+1]
			break
		}
	}

	return tag, nil
}

func findTag(lines []string) ([]int, error) {
	var tagLines []int

	for i, line := range lines {
		if strings.Contains(line, "###ahoy:tag###") {
			tagLines = append(tagLines, i)

			if len(tagLines) == 2 {
				break
			}
		}
	}
	if len(tagLines) != 2 {
		return nil, errors.New("Incorrect labeling for tag")
	}
	return tagLines, nil
}