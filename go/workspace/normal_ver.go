package main

import (
	//"bytes"
	"fmt"
	"os"
	"strconv"
	"os/exec"
	"strings"
	"time"
)

var execCommand = exec.Command

func checkReadyNode() int {
	cmd := exec.Command("kubectl", "get", "node", "--no-headers")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error executing kubectl command:", err)
		return -1
	}
	lines := strings.Split(string(output), "\n")
	cntNode := len(lines) - 1

	//cmd = exec.Command("kubectl", "get", "nodes")
	//output, err = cmd.Output()
	//if err != nil {
	//	fmt.Println("Error executing kubectl command:", err)
	//	return -1
	//}
	//lines = strings.Split(string(output), "\n")
	cntReadyNode := 0
	for _, line := range lines[1:] {
		if strings.Contains(line, " Ready ") {
			cntReadyNode++
		}
	}
	cntNotReadyNode := cntNode - cntReadyNode
	return cntNotReadyNode
}

func checkReadyPod() int {
	cmd := exec.Command("kubectl", "get", "pods")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error executing kubectl command:", err)
		return -1
	}
	lines := strings.Split(string(output), "\n")
	cntPod, cntReadyPod := 0, 0
	for _, line := range lines[1:] {
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		statusParts := strings.Split(fields[1], "/")
		if len(statusParts) != 2 {
			continue
		}
		sum, ready := statusParts[1], statusParts[0]
		cntPod += atoi(sum)
		cntReadyPod += atoi(ready)
	}
	cntNotReadyPod := cntPod - cntReadyPod
	return cntNotReadyPod
}

func getPodStatus() {
	cmd := exec.Command("kubectl", "get", "pods")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error executing kubectl command:", err)
		return -1
	}
	lines := strings.Split(string(output), "\n")

	pendingCount, creatingCount := 0, 0
	for _, line := range lines[1:] {
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}
		statusCount, status := fields[1], fields[2]
		statusParts := strings.Split(statusCount, "/")
		if len(statusParts) != 2 {
			continue
		}

		switch status {
		case "Pending":
			pendingCount += atoi(statusParts[1])
		case "ContainerCreating":
			creatingCount += atoi(statusParts[1])
		}
	}

	fmt.Printf("Pending: %d\n", pendingCount)
	fmt.Printf("ContainerCreating: %d\n", creatingCount)
}

// 하나의 문자열 인수를 받아 정수로 변환: 성공하면 변환된 정수와 nil오류, 실패하면 0과 오류
func atoi(str string) int {
	num, err := strconv.Atoi(str)
	if err != nil {
		return 10
	}
	return num
}

func main() {
	var VAR_ITER, VAR_TIME int = 0, 10
	for _, arg := range os.Args[1:] {
		parts := strings.Split(arg, "=")
		//fmt.Printf("%s\n", parts)
		if len(parts) != 2 {
			continue
		}
		key, value := parts[0], parts[1]
		switch key {
		case "iter":
			VAR_ITER = atoi(value)
		case "time":
			VAR_TIME = atoi(value)
		}
	}

	count := 0
	for {
		notReadyNode := checkReadyNode()
		notReadyPod := checkReadyPod()

		if notReadyNode == 0 && notReadyPod == 0 {
			fmt.Println("\nReady")
			break
		}

		if VAR_ITER > 0 {
			count++
			if count > VAR_ITER {
				break
			}
		}

		if notReadyPod != 0 {
			getPodStatus()
		}

		fmt.Println("Retrying ...")
		time.Sleep(time.Duration(VAR_TIME) * time.Second)
	}
}
