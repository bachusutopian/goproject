package main

import "fmt"

func matrix_multiplication(a, b, c [][]int, row, col int) {
    for k := 0; k < len(a[0]); k++ {
        c[row][col] += a[row][k] * b[k][col]
    }
}

func main() {
    // we defined the matrix
	var a [][]int
	var b [][]int
	var c [][]int
	a = [][]int{{5, 6, 7}, {1, 1, 1}}  // the matrix we chose
	b = [][]int{{2, 3}, {2, 2}, {3, 1}}  //the matrix we chose
	c = make([][]int, len(a)) //make function will create a 0 array and return a slice referencing an array 
    for i := range c {
        c[i] = make([]int, len(b[0]))
    }
/*
It uses a range loop to iterate over the elements in the variable "c", which is a matrix. 
For each iteration, the code creates a new line of integers with the length of the b[0]  and assigns it to the current element in "c" using the index variable "i".
This makes c the same size like b with elements measuring b[0] size. 
*/


    // go routines
    for i := 0; i < len(a); i++ {
        for j := 0; j < len(b[0]); j++ {
            go matrix_multiplication(a, b, c, i, j)
        }
    }

    // the classic Wait to finish
    fmt.Scanln()

    // Print the result
    fmt.Println(c)
}