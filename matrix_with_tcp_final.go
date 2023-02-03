package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
)

func matrix_multiplication(a, b, c [][]int, row int, col int) {
	print("dudu\n")
	var k int
	for k = 0; k < len(a[0]); k++ {
		print("dada\n")
		c[row][col] = a[row][k] * b[k][col]
	}
}

func openMatrix(name string) [][]int {
	// open file
	f, err := os.OpenFile(name, os.O_RDONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	// remember to close the file at the end of the program
	defer f.Close()

	// read the file line by line using scanner
	scanner := bufio.NewScanner(f)

	var mat [][]int

	for scanner.Scan() {
		var row []int

		str := scanner.Text()
		hf := strings.Split(str, " ")

		for i := 0; i < len(hf); i++ {
			j, err := strconv.Atoi(hf[i])
			if err == nil {
				row = append(row, j)
			}
		}

		mat = append(mat, row)

	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return mat
}

func handleConnection(connection net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Println("Connection accepted ", connection.RemoteAddr())

	// Send welcome message and request matrix dimensions
	connection.Write([]byte("Hello ! Welcome to the Matrix Multiplicator\n"))

	var a [][]int
	var b [][]int

	// Read names of files from client
	scanner := bufio.NewScanner(connection)

	var mat1 string
	connection.Write([]byte("Enter the name of the file for the first matrix:"))
	scanner.Scan()
	mat1 = scanner.Text()
	fmt.Printf("Matrix 1: %q\n\n", mat1)

	var mat2 string
	connection.Write([]byte("Enter the name of the file for the second matrix:"))
	scanner.Scan()
	mat2 = scanner.Text()
	fmt.Printf("Matrix 2: %q\n\n", mat2)

	a = openMatrix(mat1)
	b = openMatrix(mat2)

	dim1 := len(a)
	dim2 := len(b[0])

	var c [][]int
	c = make([][]int, dim1)
	for i := range c {
		c[i] = make([]int, dim2)
		for j := range c[i] {
			c[i][j] = 0
		}
	}

	for i := 0; i < len(a); i++ {
		for j := 0; j < len(b[i]); j++ {
			go matrix_multiplication(a, b, c, i, j)
		}
	}

	connection.Write([]byte("Multiplication done!\n"))
	for i := 0; i < len(c); i++ {
		for j := 0; j < len(c[i]); j++ {
			connection.Write([]byte(fmt.Sprintf("%d ", c[i][j])))
		}
		connection.Write([]byte("\n"))
	}

	connection.Close()
}

func main() {
	listener, error_ := net.Listen("tcp", ":9000")
	if error_ != nil {
		panic(error_)
	}
	var wg sync.WaitGroup
	for {
		connection, error_ := listener.Accept()
		if error_ != nil {
			log.Println("Accept Error", error_)
			continue
		}
		wg.Add(1)
		go func() {
			handleConnection(connection, &wg)
		}()
	}
	wg.Wait()
}
