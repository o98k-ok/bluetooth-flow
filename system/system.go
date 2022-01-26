package system

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
)

type Info struct {
	KeyWord string
	Data    map[string]interface{}
}

func NewSystem(key string) (*Info, error) {
	command := fmt.Sprintf("system_profiler %s -json > .system.log", key)
	d, err := exec.Command("bash", "-c", command).CombinedOutput()
	if err != nil {
		return nil, err
	}

	// fmt.Printf("%s", string(d))
	d, err = os.ReadFile(".system.log")
	if err != nil {
		return nil, err
	}

	var res map[string]interface{}
	if err = json.Unmarshal(d, &res); err != nil {
		return nil, err
	}
	return &Info{key, res}, nil
}
