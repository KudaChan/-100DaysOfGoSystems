package main

import "fmt"

func main() {
	arr := []int{1, 2, 3, 4, 5}

	i, j := 0, len(arr)-1
	for i < j {
		arr[i], arr[j] = arr[j], arr[i]
		i++
		j--
	}
	fmt.Println(arr)
}
