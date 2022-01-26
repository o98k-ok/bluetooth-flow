package system

import "testing"

func TestGetAllBlueTooth(t *testing.T) {
	err := InitBlueTooth()
	if err != nil {
		t.Fail()
	}

	t.Log(GetAllBlueTooth())
}
