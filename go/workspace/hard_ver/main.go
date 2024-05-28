package main

import (
    "fmt"
    "log"
    "context"
    "filepath"

    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/tools/clientcmd"
    "k8s.io/client-go/util/homedir"
)

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

    fmt.Println("<Nodes>")
    for _, node := range nodes.Items {
        fmt.Printf("Name: %s\n", node.Name)
    }


	
}
