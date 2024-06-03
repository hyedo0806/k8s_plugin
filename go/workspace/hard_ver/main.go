package main

import (
    //"os"
    "time"
    "flag"
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

    //wordPtr := flag.String("iter", "time")
    iterPtr := flag.Int("iter", 0, "default: 10iterations")
    timePtr := flag.Int("time", 10, "default: 10sec")
    flag.Parse()
    //fmt.Println("iter : %d", *iterPtr)
    //fmt.Println("time : %d", *timePtr)
    
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

    for {    
        // 노드 정보 가져오기
        nodes, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
        if err != nil {
            log.Fatalf("Error listing nodes: %v", err)
        }

        cntNode := len(nodes.Items)
        cntReadyNode := 0
        //fmt.Println("<Nodes>")

        for _, node := range nodes.Items {
            //fmt.Printf("Name: %s\n", node.Name)

            data := node.Status.Conditions
            //fmt.Println(data[0].Status)
            for _, condition := range data {

                if condition.Status == "True" && condition.Type == "Ready" {
                    cntReadyNode += 1
                }
                //fmt.Printf("%d Type: %s, Status: %s\n",i, condition.Type, condition.Status)
            }

        }
        fmt.Printf("Total Node %d\n", cntNode)
        fmt.Printf("Ready Node %d\n", cntReadyNode)

        // 파드 정보 가져오기 (디폴트 네임스페이스)
        pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
        if err != nil {
            log.Fatalf("Error listing pods: %v", err)
        }

        cntPod := len(pods.Items)
        cntReadyPod := 0
        //fmt.Println("<Pods>")
        for _, pod := range pods.Items {
            //fmt.Printf("Namespace: %s, Name: %s\n", pod.Namespace, pod.Name)

            data := pod.Status.Conditions
            //fmt.Println(data[0].Status)
            for _, condition := range data {

                    if condition.Status == "True" && condition.Type == "Ready" {

                            cntReadyPod += 1
                    }
                    //fmt.Printf("%d Type: %s, Status: %s\n",i, condition.Type, condition.Status)
            }
        }
        fmt.Printf("Total Pod %d\n", cntPod)
        fmt.Printf("Ready Pod %d\n", cntReadyPod)

        if cntNode != cntReadyNode || cntPod != cntReadyPod {
            fmt.Println("Retrying ...")
		    time.Sleep(10)

        } else {
            fmt.Println("Pass")
            break
        }
    }
}
