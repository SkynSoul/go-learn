package utils

import (
	"testing"
)

func TestGetExecuteFilePath(t *testing.T) {
	t.Log(GetExecuteFilePath())
}


func TestGetWorkingPath(t *testing.T) {
	t.Log(GetWorkingPath())
}