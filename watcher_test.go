package main

import (
	"github.com/mitchellh/go-homedir"
	"testing"
)

func TestListSubDir(t *testing.T) {
	dir, _ := homedir.Expand("~/Downloads")
	t.Log(ListSubDir(dir))
}
