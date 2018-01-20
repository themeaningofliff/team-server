package main

import "fmt"

func add(x, y int) (result int) {
	result = x + y
	return

}
func main() {
	fmt.Printf("hello, world. Did you know 1+3=%d\n", add(1, 3))
}
