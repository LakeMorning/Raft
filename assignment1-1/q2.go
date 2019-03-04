package cos418_hw1_1

import (
	"bufio"
	"io"
	"log"
	"os"
	"strconv"
	"sync"
)

// Sum numbers from channel `nums` and output sum to `out`.
// You should only output to `out` once.
// Do NOT modify function signature.
func sumWorker(nums chan int, out chan int) {
	// TODO: implement me
	// HINT: use for loop over `nums`
	res := 0
	for num := range nums {
		res += num
	}
	out <- res
}

// Read integers from the file `fileName` and return sum of all values.
// This function must launch `num` go routines running
// `sumWorker` to find the sum of the values concurrently.
// You should use `checkError` to handle potential errors.
// Do NOT modify function signature.
func sum(num int, fileName string) int {
	// TODO: implement me
	// HINT: use `readInts` and `sumWorkers`
	// HINT: used buffered channels for splitting numbers between workers
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	elems, err := readInts(file)
	if err != nil {
		log.Fatal(err)
	}

	nums := make(chan int)
	out := make(chan int, num)

	var wg sync.WaitGroup
	for i := 0; i < num; i++ {
		wg.Add(1)
		go func() {
			sumWorker(nums, out)
			wg.Done()
		}()
	}

	go func() {
		wg.Add(1)
		for _, elem := range elems {
			nums <- elem
		}
		wg.Done()
		close(nums)
	}()

	wg.Wait()
	close(out)

	res := 0
	for x := range out {
		res += x
	}
	return res
}

// Read a list of integers separated by whitespace from `r`.
// Return the integers successfully read with no error, or
// an empty slice of integers and the error that occurred.
// Do NOT modify this function.
func readInts(r io.Reader) ([]int, error) {
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanWords)
	var elems []int
	for scanner.Scan() {
		val, err := strconv.Atoi(scanner.Text())
		if err != nil {
			return elems, err
		}
		elems = append(elems, val)
	}
	return elems, nil
}
