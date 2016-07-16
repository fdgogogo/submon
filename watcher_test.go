package main

import (
	"testing"
	"github.com/mitchellh/go-homedir"
)

func TestListSubDir(t *testing.T) {
	dir, _ := homedir.Expand("~/Downloads")
	t.Log(ListSubDir(dir))
}
