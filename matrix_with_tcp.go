package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

func matrix_multiplication(a, b, c [][]int, row, col int) {
	var k int
	for k = 0; k < len(a[0]); k++ {
		c[row][col] += a[row][k] * b[k][col]
	}
}

func main() {
	listener, err := net.Listen("tcp", ":9000")
	if err != nil {
		panic(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Accept Error", err)
			continue
		}

		log.Println("Accepted ", conn.RemoteAddr())

		// Send welcome message and request matrix dimensions
		conn.Write([]byte("Hello ! Welcome to the Matrix Multiplicator\n"))
		conn.Write([]byte("Please enter the dimensions of the first matrix (rows columns):\n>"))

		// Read matrix dimensions from client
		s := bufio.NewScanner(conn)
		var a_rows, a_cols, b_rows, b_cols int
		if s.Scan() {
			input := s.Text()
			n, err := fmt.Sscanf(input, "%d %d", &a_rows, &a_cols)
			if n != 2 || err != nil {
				conn.Write([]byte("This is an invalid input, tap the right rows and columns\n"))
				conn.Close()
				continue
			}
		}
		conn.Write([]byte("Please, enter the first matrix:\n>"))
		a := make([][]int, a_rows)
		for i := range a {
			a[i] = make([]int, a_cols)
		}
		for i := 0; i < a_rows; i++ {
			for j := 0; j < a_cols; j++ {
				if !s.Scan() {
					conn.Write([]byte("Invalid input, please enter the matrix elements separated by spaces\n"))
					conn.Close()
					continue
				}
				_, err := fmt.Sscanf(s.Text(), "%d", &a[i][j])
				if err != nil {
					conn.Write([]byte("Invalid input, please enter integer values for the matrix elements\n"))
					conn.Close()
					continue
				}
			}
		}
		conn.Write([]byte("Please enter the dimensions of the second matrix (rows columns):\n>"))
		if s.Scan() {
			input := s.Text()
			n, err := fmt.Sscanf(input, "%d %d", &b_rows, &b_cols)
			if n != 2 || err != nil {
				conn.Write([]byte("Invalid input, please enter the dimensions in the format 'rows columns'\n"))
				conn.Close()
				continue
			}
		}
		if a_cols != b_rows {
			conn.Write([]byte("Invalid matrix dimensions, the number of columns of the first matrix must be equal to the number of rows of the second matrix\n"))
			conn.Close()
			continue
		}
		conn.Write([]byte("Please enter the second matrix:\n>"))
		b := make([][]int, b_rows)
		for i := range b {
			b[i] = make([]int, b_cols)
		}
		for i := 0; i < b_rows; i++ {
			for j := 0; j < b_cols; j++ {
				if !s.Scan() {
					conn.Write([]byte("Invalid input, please enter the matrix elements separated by spaces\n"))
					conn.Close()
					continue
				}
				_, err := fmt.Sscanf(s.Text(), "%d", &b[i][j])
				if err != nil {
					conn.Write([]byte("Invalid input, please enter integer values for the matrix elements\n"))
					conn.Close()
					continue
				}
			}
		}

		c := make([][]int, a_rows)
		for i := range c {
			c[i] = make([]int, b_cols)
		}

		// matrix multiplication
		for i := 0; i < a_rows; i++ {
			for j := 0; j < b_cols; j++ {
				go matrix_multiplication(a, b, c, i, j)
			}
		}

		// Wait for goroutines to finish
		fmt.Scanln()

		// Send result matrix to client
		for i := 0; i < a_rows; i++ {
			for j := 0; j < b_cols; j++ {
				conn.Write([]byte(fmt.Sprintf("%d ", c[i][j])))
			}
			conn.Write([]byte("\n"))
		}

		// Close connection
		conn.Close()
	}
}
