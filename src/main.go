package main

import "fmt"

func gcd(numbers ...int) (res int) {
	res = numbers[0]
	for i := 1; i < len(numbers); i++ {
		a, b := numbers[i], res%numbers[i]
		for b != 0 {
			a, b = b, a%b
		}
		res = a
	}
	return
}

func main() {
	a, b, c := 10, 20, 25
	fmt.Println(gcd(a, b, c))
}
