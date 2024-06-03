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
	//	fmt.Printf("Name: %s\n", node.Name)
    	//fmt.Printf("Labels: %v\n", node.Labels)
    	//fmt.Printf("Annotations: %v\n", node.Annotations)
    	//fmt.Printf("Addresses: %v\n", node.Status.Addresses)
    	fmt.Printf("Conditions: %v\n", node.Status.Conditions)
    	//fmt.Printf("Capacity: %v\n", node.Status.Capacity)
    	//fmt.Printf("Allocatable: %v\n", node.Status.Allocatable)
    	//fmt.Printf("DaemonEndpoints: %v\n", node.Status.DaemonEndpoints)
    	//fmt.Printf("NodeInfo: %v\n", node.Status.NodeInfo)
    	//fmt.Printf("PodCIDR: %v\n", node.Spec.PodCIDR)
    	//fmt.Printf("ProviderID: %v\n", node.Spec.ProviderID)
    	//fmt.Printf("Unschedulable: %v\n", node.Spec.Unschedulable)
    }
    
    // 파드 정보 가져오기 (디폴트 네임스페이스)
    pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
    if err != nil {
        log.Fatalf("Error listing pods: %v", err)
    }

    fmt.Println("<Pods>")
    for _, pod := range pods.Items {
        fmt.Printf("Namespace: %s, Name: %s\n", pod.Namespace, pod.Name)
	labels := pod.Labels
		for k, v := range labels {
			fmt.Printf("Pod %s Namespace %s Label %s:%s ", pod.Name, pod.Namespace, k, v)
			fmt.Printf("Contidions :%s\n", pod.Status)
		}
	}


}
