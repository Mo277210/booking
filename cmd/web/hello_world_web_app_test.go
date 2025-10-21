package main

import "testing"

func TestRun(t *testing.T) {
	db,err := run()
	if err != nil {
		t.Errorf("run() returned an error: %v", err)
	}
	defer db.SQL.Close()

}
