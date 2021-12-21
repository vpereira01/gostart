// Hello Word in Go by Vivek Gite
package main

// Import OS and fmt packages
import (
	"fmt"
	"os"

	"github.com/NOX73/go-neural"
	"github.com/NOX73/go-neural/engine"
)

// Let us start
func main() {
	fmt.Println("Hello, world!")                          // Print simple text on screen
	fmt.Println(os.Getenv("USER"), ", Let's be friends!") // Read Linux $USER environment variable

	network := neural.NewNetwork(2, []int{2, 2})
	engine := engine.New(network)
	engine.Start()

	engine.Learn([]float64{1, 2}, []float64{3, 3}, 0.1)
	engine.Calculate([]float64{1, 2})
	fmt.Printf("engine.Dump(): %v\n", engine.Dump())

	fmt.Printf("engine.Calculate([]float64{1, 2}): %v\n", engine.Calculate([]float64{1, 2}))
}
