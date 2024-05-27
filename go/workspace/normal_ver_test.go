// main_test.go
package main

import (
	"fmt"
	"os"
	"os/exec"
	//"strings"
	"testing"
)

func TestAtoi(t *testing.T) {
	if atoi("42") != 42 {
		t.Errorf("Expected 42, got %d", atoi("42"))
	}
	if atoi("0") != 0 {
		t.Errorf("Expected 0, got %d", atoi("0"))
	}
	if atoi("abc") != 0 {
		t.Errorf("Expected 0 for invalid input, got %d", atoi("abc"))
	}
}

func mockCommand(name string, arg ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess", "--", name}
	cs = append(cs, arg...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}

func TestCheckReadyNode(t *testing.T) {
	execCommand = mockCommand
	defer func() { execCommand = exec.Command }()

	expected := 1
	actual := checkReadyNode()
	if actual != expected {
		t.Errorf("Expected %d, got %d", expected, actual)
	}
}

func TestCheckReadyPod(t *testing.T) {
	execCommand = mockCommand
	defer func() { execCommand = exec.Command }()

	expected := 2
	actual := checkReadyPod()
	if actual != expected {
		t.Errorf("Expected %d, got %d", expected, actual)
	}
}

func TestHelperProcess(*testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	cmd := os.Args[3]
	switch cmd {
	case "kubectl get node --no-headers":
		fmt.Fprint(os.Stdout, "node1\nnode2\n")
	case "kubectl get nodes":
		fmt.Fprint(os.Stdout, "NAME STATUS\nnode1 Ready\nnode2 NotReady\n")
	case "kubectl get pods":
		fmt.Fprint(os.Stdout, "NAME READY STATUS\npod1 0/1 Pending\npod2 1/1 Running\npod3 0/1 ContainerCreating\n")
	}
	os.Exit(0)
}
