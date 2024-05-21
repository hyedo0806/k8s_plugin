package main

import (
    "fmt"
    "os/exec"
    "strconv"
    "strings"
)

func main() {
    fmt.Printf("\n<Nodes>\n")
    // Ready 상태의 노드 수 가져오기
    cntReadyNode, err := exec.Command("sh", "-c", "kubectl get nodes | grep -c ' Ready'").Output()
    if err != nil {
        fmt.Printf("Error getting ready nodes: %v\n", err)
        return
    }
    fmt.Printf("# ready_node: %s\n", strings.TrimSpace(string(cntReadyNode)))

    // NotReady 상태의 노드 수 가져오기
    cntNotReadyNode, err := exec.Command("sh", "-c", "kubectl get nodes | grep -c ' NotReady'").Output()
    if err != nil {
        fmt.Printf("Error getting not ready nodes: %v\n", err)
        return
    }
    fmt.Printf("# notready_node: %s\n", strings.TrimSpace(string(cntNotReadyNode)))

    fmt.Printf("\n<Pods>\n")
    // 모든 파드의 수 가져오기
    cntPod, err := exec.Command("sh", "-c", "kubectl get pods | awk 'NR>1 {split($2, a, \"/\"); sum += a[2]} END {print sum}'").Output()
    if err != nil {
        fmt.Printf("Error getting pods count: %v\n", err)
        return
    }
    cntPodStr := strings.TrimSpace(string(cntPod))
    fmt.Printf("# pod: %s\n", cntPodStr)

    // Ready 상태의 파드 수 가져오기
    cntReadyPod, err := exec.Command("sh", "-c", "kubectl get pods | awk 'NR>1 {split($2, a, \"/\"); sum += a[1]} END {print sum}'").Output()
    if err != nil {
        fmt.Printf("Error getting ready pods count: %v\n", err)
        return
    }
    cntReadyPodStr := strings.TrimSpace(string(cntReadyPod))
    fmt.Printf("# ready_pod: %s\n", cntReadyPodStr)

    // NotReady 상태의 파드 수 계산
    cntNotReadyPod := fmt.Sprintf("%d", toInt(cntPodStr)-toInt(cntReadyPodStr))
    fmt.Printf("# notready_pod: %s\n", cntNotReadyPod)

    // Pending, Failed, Unknown 상태의 파드 수 계산
    if toInt(cntNotReadyPod) != 0 {
        podsInfo, err := exec.Command("sh", "-c", "kubectl get pods").Output()
        if err != nil {
            fmt.Printf("Error getting pods info: %v\n", err)
            return
        }
        podsInfoStr := string(podsInfo)	
	//fmt.Printf("podInfo : %s\n", podsInfoStr)
	//terminating := strings.Count(podsInfoStr, "Terminating")
        //pending := strings.Count(podsInfoStr, "Pending")
        //failed := strings.Count(podsInfoStr, "Failed")
        //unknown := strings.Count(podsInfoStr, "Unknown")

        //fmt.Printf("Terminating pod: %d\n", terminating)
	//fmt.Printf("Pending pod: %d\n", pending)
        //fmt.Printf("Failed pod: %d\n", failed)
        //fmt.Printf("Unknown pod: %d\n", unknown)
	// 포드 정보를 줄 단위로 분할
	lines := strings.Split(podsInfoStr, "\n")

	// 첫 번째 줄은 헤더이므로 무시
	lines = lines[1:]

	pendingCount := 0
	terminatingCount := 0

	for _, line := range lines {

		if line == "" {
			continue
		}
		//fmt.Printf("%d. %s :  ", i, line)
		// 각 줄을 공백으로 분할하여 필드를 가져옴
		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}

		statusCount := fields[1]
		status := fields[2]

		// 슬래시 기호로 READY 상태를 분할
		statusParts := strings.Split(statusCount, "/")
		if len(statusParts) != 2 {
			continue
		}

		// 상태별로 카운트 증가
		switch status {
		case "Pending":
			pendingCount+=toInt(statusParts[1])
		case "Terminating":
			terminatingCount+=toInt(statusParts[1])
		}
	}
	fmt.Printf("	Pending : %d\n", pendingCount)
	//fmt.Printf("    Terminatinging : %d\n", terminatingCount)
    }


    // 클러스터 상태 체크
    cntNode := toInt(strings.TrimSpace(string(cntReadyNode))) + toInt(strings.TrimSpace(string(cntNotReadyNode)))
    if cntNode == toInt(strings.TrimSpace(string(cntReadyNode))) && toInt(cntPodStr) == toInt(cntReadyPodStr) {
        fmt.Println("\npass!!")
    } else {
        fmt.Println("\nerror!!")
    }
}

func toInt(str string) int {
    val, err := strconv.Atoi(strings.TrimSpace(str))
    if err != nil {
        fmt.Printf("Error converting string to int: %v\n", err)
        return 0
    }
    return val
}

