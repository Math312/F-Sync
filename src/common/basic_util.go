package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func GetOsType() string {
	return runtime.GOOS
}
func ParseUrlFragmentParameters(urlStr string) (map[string]string, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	fragment := u.Fragment
	parameters := strings.Split(fragment, "&")
	parameterMap := make(map[string]string, 5)
	for _, str := range parameters {
		data := strings.Split(str, "=")
		parameterMap[data[0]] = data[1]
	}
	return parameterMap, nil
}

func ReadJson[T any](config []byte, result *T) error {
	if err := json.Unmarshal(config, result); err != nil {
		return errors.New("Json Config Parse Error ")
	} else {
		return nil
	}
}

func WriteJson[T any](config *T) ([]byte, error) {
	if result, err := json.Marshal(*config); err != nil {
		return nil, err
	} else {
		return result, nil
	}
}

func FileIsExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		if os.IsNotExist(err) {
			return false
		}
		fmt.Println(err)
		return false
	}
	return true
}

var commands = map[string]string{
	"windows": "cmd /c start",
	"darwin":  "open",
	"linux":   "xdg-open",
}

// Open calls the OS default program for uri
func OpenUri(uri string) error {
	run, ok := commands[runtime.GOOS]
	if !ok {
		return fmt.Errorf("don't know how to open things on %s platform", runtime.GOOS)
	}

	cmd := exec.Command(run, uri)
	return cmd.Start()
}
