package spec

import (
	"os"
	"regexp"
	"strings"

	"github.com/BurntSushi/toml"
)

var argRx = regexp.MustCompile(`(?:[^\s"]+|"[^"]*")+`)

type LanguageSpec struct {
	Cmd      string `json:"cmd"      toml:"cmd"`
	FileName string `json:"fileName" toml:"fileName"`
	Language string `json:"language" toml:"language"`
}

func (ls LanguageSpec) String() string {
	return ls.Cmd + " " + ls.FileName
}

func (ls *LanguageSpec) GetCommandWithArgs() []string {
	return split(ls.Cmd)
}

type Spec map[string]LanguageSpec

func New(filePath string) (*Spec, error) {
	// Read file
	file, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var spec Spec
	_, err = toml.Decode(string(file), &spec)
	if err != nil {
		return nil, err
	}
	return &spec, nil
}

func (s *Spec) Get(language string) (*LanguageSpec, error) {
	langSpec, ok := (*s)[language]
	if !ok {
		return nil, os.ErrNotExist
	}
	return &langSpec, nil
}

func split(v string) (res []string) {
	res = argRx.FindAllString(v, -1)
	for i, v := range res {
		res[i] = strings.ReplaceAll(v, "\"", "")
	}
	return
}
