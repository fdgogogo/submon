package main

import "testing"

func TestWorkerGroup_Run(t *testing.T) {
	t.Log(NewTaskQueue().Run())
}
