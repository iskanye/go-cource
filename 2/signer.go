package main

import (
	"sort"
	"strconv"
	"strings"
	"sync"
)

func ExecutePipeline(jobs ...job) {
	in := make(chan interface{})
	var wg sync.WaitGroup

	for _, job := range jobs {
		out := make(chan interface{})
		wg.Add(1)

		go func() {
			defer wg.Done()
			job(in, out)
			close(out)
		}()

		in = out
	}

	wg.Wait()
}

func SingleHash(in, out chan interface{}) {
	var wg sync.WaitGroup

	for i := range in {
		data := i.(string)
		wg.Add(1)

		go func() {
			defer wg.Done()
			out <- (<-Crc32Channel(data)) + "~" + (<-Crc32Channel(DataSignerMd5(data)))
		}()
	}

	wg.Wait()
}

func MultiHash(in, out chan interface{}) {
	var wg sync.WaitGroup

	for i := range in {
		data := i.(string)
		wg.Add(1)

		go func() {
			defer wg.Done()

			channels := make([]chan string, 0, 6)

			for j := 0; j < 6; j++ {
				channels = append(channels, Crc32Channel(strconv.Itoa(j)+data))
			}

			output := make([]string, 0, 6)

			for _, ch := range channels {
				output = append(output, <-ch)
			}

			out <- strings.Join(output, "")
		}()
	}

	wg.Wait()
}

func CombineResults(in, out chan interface{}) {
	result := make([]string, 0)

	for i := range in {
		result = append(result, i.(string))
	}

	sort.Strings(result)
	out <- strings.Join(result, "_")
}

func Crc32Channel(str string) chan string {
	out := make(chan string, 1)

	go func() {
		out <- DataSignerCrc32(str)
		close(out)
	}()

	return out
}
