package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"
)

func main() {
	start := time.Now()

	fd, err := os.Open("input.txt")
	if err != nil {
		fmt.Println(err)
	}
	defer fd.Close()

	scanner := bufio.NewScanner(fd)

	inputs := make([]chan string, runtime.NumCPU())
	for i := 0; i < runtime.NumCPU(); i++ {
		inputs[i] = make(chan string)
	}

	output := make(chan string, runtime.NumCPU())
	done := make(chan struct{})

	wg := sync.WaitGroup{}
	wg.Add(runtime.NumCPU())

	for i := 0; i < runtime.NumCPU(); i++ {
		go pigLatinEncoderWorker(inputs[i], output, &wg)
	}

	go outputWriter(output, done)

	counter := 0
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			fmt.Println("Received line from scanner..")
			inputs[counter%runtime.NumCPU()] <- line
			counter++
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}

	//all paragraphs are read..close the inputs
	for _, c := range inputs {
		close(c)
	}

	wg.Wait()     //wait for workers to finish
	close(output) //workers exited..close the output
	<-done        //outputWriter finished

	fmt.Println(time.Since(start))
}

func outputWriter(out chan string, done chan struct{}) {
	fd, err := os.Create("output.txt")
	if err != nil {
		fmt.Println(err)
	}
	defer fd.Close()

	for line := range out {
		fmt.Println("Received line from output..")
		fmt.Fprintln(fd, line)
		if err != nil {
			fmt.Println(err)
		}
	}
	done <- struct{}{}
}

func pigLatinEncoderWorker(in chan string, out chan string, wg *sync.WaitGroup) {
	for line := range in {
		//convert to pig latin
		res := ""
		temp := ""

		for _, ch := range line {
			if ch == ' ' || ch == '.' || ch == ',' || ch == ';' || ch == '!' || ch == ':' {
				if temp == "" {
					res += string(ch)
					continue
				}

				res += pigLatinEncode(temp) + string(ch)
				temp = ""
			} else {
				temp += string(ch)
			}
		}

		if temp != "" {
			res += pigLatinEncode(temp)
		}

		out <- res
	}
	wg.Done()
}

func pigLatinEncode(word string) string {

	if word[0] == 'a' || word[0] == 'A' ||
		word[0] == 'e' || word[0] == 'E' ||
		word[0] == 'i' || word[0] == 'I' ||
		word[0] == 'o' || word[0] == 'O' ||
		word[0] == 'u' || word[0] == 'U' {
		return word + "ay"
	} else {
		return word[1:] + string(word[0]) + "ay"
	}
}
