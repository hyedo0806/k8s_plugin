package main

import (
	"fmt"
	"go/workspace/usepkg/printfpkg"

	"github.com/guptarohit/asciigraph"
	"github.com/tuckersGo/musthaveGo/ch16/expkg"
)

func main() {
    	printfpkg.PrintHello()
    	expkg.PrintSample()
	
			data := []float64{3, 4, 5, 6, 9,7,8,5,10,2,7,5,6}
			graph := asciigraph.Plot(data)
			fmt.Println(graph)
    	//fmt.Printf("Hello, World\n")
}
