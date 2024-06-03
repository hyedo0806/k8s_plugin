package main

import (
	"testing"
)

func TestAtoi(t *testing.T) {
	if atoi("42") != 42 {
		t.Errorf("Expected 42, got %d", atoi("42"))
	}
	if atoi("0") != 0 {
		t.Errorf("Expected 0, got %d", atoi("0"))
	}
}

func TestCheckReadyNode(t *testing.T) {
	expected := 1 // Adjust this value based on your actual cluster status
	actual := checkReadyNode()
	if actual != expected {
		t.Errorf("Expected %d, got %d", expected, actual)
	}
}

func TestCheckReadyPod(t *testing.T) {
	expected := 2 // Adjust this value based on your actual cluster status
	actual := checkReadyPod()
	if actual != expected {
		t.Errorf("Expected %d, got %d", expected, actual)
	}
}
