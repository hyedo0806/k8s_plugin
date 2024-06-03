package main

import (
    "fmt"
    "log"
    "context"
    "path/filepath"

    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/tools/clientcmd"
    "k8s.io/client-go/util/homedir"
)

type Condition struct {
	Type               string
	Status             string
	LastHeartbeatTime  string
	LastTransitionTime string
	Reason             string
	Message            string
}

func parseConditions(input string) []Condition {
	// Regular expression to match each condition block
	re := regexp.MustCompile(`{([^}]*)}`)
	matches := re.FindAllStringSubmatch(input, -1)

	var conditions []Condition

	for _, match := range matches {
		// Splitting each condition block into fields
		fields := strings.Fields(match[1])
		if len(fields) >= 6 {
			condition := Condition{
				Type:               fields[0],
				Status:             fields[1],
				LastHeartbeatTime:  fields[2] + " " + fields[3] + " " + fields[4] + " " + fields[5],
				LastTransitionTime: fields[6] + " " + fields[7] + " " + fields[8] + " " + fields[9],
				Reason:             fields[10],
				Message:            strings.Join(fields[11:], " "),
			}
			conditions = append(conditions, condition)
		}
	}

	return conditions
}

func main() {
     // kubeconfig 파일 경로 설정
     var kubeconfig string
     if home := homedir.HomeDir(); home != "" {
         kubeconfig = filepath.Join(home, ".kube", "config")
     } else {
         log.Fatal("Unable to find kubeconfig file")
     }

    // 클라이언트 설정 및 초기화
    config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
    if err != nil {
        log.Fatalf("Error building kubeconfig: %v", err)
    }

    clientset, err := kubernetes.NewForConfig(config)
    if err != nil {
        log.Fatalf("Error creating Kubernetes client: %v", err)
    }


    // 노드 정보 가져오기
    nodes, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
    if err != nil {
        log.Fatalf("Error listing nodes: %v", err)
    }

    cntNode := len(nodes.Items)
    cntReadyNode := 0
    fmt.Println("<Nodes> %d\n",cntNode)
    for _, node := range nodes.Items {
        fmt.Printf("Name: %s\n", node.Name)

        input = strings.TrimPrefix(node.Status.Conditions, "Conditions: ")

	    conditions := parseConditions(input)

        for _, condition := range conditions {
            fmt.Printf("Type: %s", condition.Type)
            fmt.Printf("  Status: %s\n", condition.Status)
            if condition.Status == "Ready" {
                cntReadyNode += 1
            }
            //fmt.Printf("LastHeartbeatTime: %s\n", condition.LastHeartbeatTime)
            //fmt.Printf("LastTransitionTime: %s\n", condition.LastTransitionTime)
            //fmt.Printf("Reason: %s\n", condition.Reason)
            //fmt.Printf("Message: %s\n", condition.Message)
            fmt.Println()
        }
    		
    }
    fmt.Printf("ReadyNode %d\n", cntReadyNode)

    // 파드 정보 가져오기 (디폴트 네임스페이스)
    pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
    if err != nil {
        log.Fatalf("Error listing pods: %v", err)
    }

    cntPod := len(pods.Items)
    cntReadyPod := 0
    fmt.Println("<Pods> %d\n", cntPod)
    for _, pod := range pods.Items {

        fmt.Printf("Namespace: %s, Name: %s\n", pod.Namespace, pod.Name)
        
        input = strings.TrimPrefix(pod.Status.Conditions, "Conditions: ")

	    conditions := parseConditions(input)

        for _, condition := range conditions {
            fmt.Printf("Type: %s", condition.Type)
            fmt.Printf("  Status: %s\n", condition.Status)
            if condition.Status == "Ready" {
                cntReadyPod += 1
            }
            //fmt.Printf("LastHeartbeatTime: %s\n", condition.LastHeartbeatTime)
            //fmt.Printf("LastTransitionTime: %s\n", condition.LastTransitionTime)
            //fmt.Printf("Reason: %s\n", condition.Reason)
            //fmt.Printf("Message: %s\n", condition.Message)
            fmt.Println()
        }
   
	}
    fmt.Printf("ReadyPod %d\n", cntReadyPod)

}
