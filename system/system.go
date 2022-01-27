package system

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

type Info struct {
	KeyWord string
	Data    map[string]interface{}
}

func NewSystem(key string) (*Info, error) {
	command := fmt.Sprintf("system_profiler %s -json 2> /dev/null", key)
	d, err := exec.Command("bash", "-c", command).CombinedOutput()
	if err != nil {
		return nil, err
	}

	var res map[string]interface{}
	if err = json.Unmarshal(d, &res); err != nil {
		return nil, err
	}
	return &Info{key, res}, nil
}
