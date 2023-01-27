package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"sync"
)

func matrix_multiplication(a, b, c [][]int, row, col int, mutex_as_parameter *sync.Mutex) {
	var k int
	for k = 0; k < len(a[0]); k++ {
		mutex_as_parameter.Lock()
		c[row][col] += a[row][k] * b[k][col]
		mutex_as_parameter.Unlock()
	}
}

func handleConnection(connection net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Println("Connection accepted ", connection.RemoteAddr())

	// Send welcome message and request matrix dimensions
	connection.Write([]byte("Hello ! Welcome to the Matrix Multiplicator\n"))
	connection.Write([]byte("Please enter the dimensions of the first matrix (rows columns format):\n_"))

	// Read matrix dimensions from client
	scan := bufio.NewScanner(connection)
	var a_rows, a_cols, b_rows, b_cols int
	if scan.Scan() {
		input := scan.Text()
		n_input, error_ := fmt.Sscanf(input, "%d %d", &a_rows, &a_cols)
		if n_input != 2 || error_ != nil {
			connection.Write([]byte("This is an invalid input, tap the right rows and columns\n"))
			connection.Close()
			return
		}
	}
	connection.Write([]byte("Please, enter the first matrix(each value on a matrix separated by an ENTER):\n_"))
	a := make([][]int, a_rows)
	for i := range a {
		a[i] = make([]int, a_cols)
	}
	for i := 0; i < a_rows; i++ {
		for j := 0; j < a_cols; j++ {
			if !scan.Scan() {
				connection.Write([]byte("Invalid input, please enter the matrix elements separated by spaces\n"))
				connection.Close()
				return
			}
			_, err := fmt.Sscanf(scan.Text(), "%d", &a[i][j])
			if err != nil {
				connection.Write([]byte("Invalid input, please enter integer values for the matrix elements\n"))
				connection.Close()
				return
			}
		}
	}
	connection.Write([]byte("Please enter the dimensions of the second matrix (rows columns format):\n>"))
	if scan.Scan() {
		input := scan.Text()
		n_input, error_ := fmt.Sscanf(input, "%d %d", &b_rows, &b_cols)
		if n_input != 2 || error_ != nil {
			connection.Write([]byte("Invalid input, please enter the dimensions in the format 'rows columns'\n"))
			connection.Close()
			return
		}
	}
	if a_cols != b_rows {
		connection.Write([]byte("Invalid matrix dimensions, the number of columns of the first matrix must be equal to the number of rows of the second matrix\n"))
		connection.Close()
		return
	}
	connection.Write([]byte("Please enter the second matrix(each value on a matrix separated by an ENTER):\n"))
	b := make([][]int, b_rows)
	for i := range b {
		b[i] = make([]int, b_cols)
	}
	for i := 0; i < b_rows; i++ {
		for j := 0; j < b_cols; j++ {
			if !scan.Scan() {
				connection.Write([]byte("Invalid input, please enter the matrix elements separated by spaces\n"))
				connection.Close()
				return
			}
			_, error_ := fmt.Sscanf(scan.Text(), "%d", &b[i][j])
			if error_ != nil {
				connection.Write([]byte("Invalid input, please enter integer values for the matrix elements\n"))
				connection.Close()
				return
			}
		}
	}
	c := make([][]int, a_rows)
	for i := range c {
		c[i] = make([]int, b_cols)
	}
	var mutex_as_parameter sync.Mutex
	for i := 0; i < a_rows; i++ {
		for j := 0; j < b_cols; j++ {
			go matrix_multiplication(a, b, c, i, j, &mutex_as_parameter)
		}
	}
	connection.Write([]byte("Multiplication done!\n"))
	for i := 0; i < a_rows; i++ {
		for j := 0; j < b_cols; j++ {
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
