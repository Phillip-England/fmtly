package main

import (
	"os"
	"os/exec"
	"testing"
)

func TestMain(t *testing.T) {
	cmd := exec.Command("go", "build", "main.go")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(), "PATH="+os.Getenv("PATH"))

	err := cmd.Run()
	if err != nil {
		t.Fatalf("Error: %s", err)
	}

	cmd = exec.Command("./main", "build", "./test/test_components", "./output.go", "main")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		t.Fatalf("Error: %s", err)
	}
}
