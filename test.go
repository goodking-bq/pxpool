package main

import "./scanner"

func main() {
	var ports = []int{80, 5000}
	scanner.ScanIP("127.0.0.1", ports)
}
