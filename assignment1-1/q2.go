package cos418_hw1_1

import (
	"bufio"
	"io"
	"os"
	"strconv"
	"sync"
)

// Sum numbers from channel `nums` and output sum to `out`.
// You should only output to `out` once.
// Do NOT modify function signature.
func sumWorker(wg *sync.WaitGroup, nums chan int, out chan int) {
	defer wg.Done()
	res := 0

	for val := range nums {
		res += val
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
	checkError(err)
	defer file.Close()

	nums, err := readInts(file)
	checkError(err)

	channels := make([]chan int, num)
	for i := range channels {
		channels[i] = make(chan int, len(nums)/num+1)
	}

	out := make(chan int, num)
	var wg sync.WaitGroup
	wg.Add(num)

	for i := 0; i < num; i++ {
		go sumWorker(&wg, channels[i], out)
	}

	for i, number := range nums {
		channels[i%num] <- number
	}

	for i := range channels {
		close(channels[i])
	}

	wg.Wait()
	close(out)

	totalSum := 0
	for partialSum := range out {
		totalSum += partialSum
	}

	return totalSum
}



// Read a list of integers from `r`.
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
