package main

import (
	"fmt"
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
		jobFunc := JobFunc(job, &wg)

		wg.Add(1)
		go jobFunc(in, out)
		in = out
	}

	wg.Wait()
}

func JobFunc(jobFunc job, wg *sync.WaitGroup) job {
	return func(in, out chan interface{}) {
		defer wg.Done()
		jobFunc(in, out)
		close(out)
	}
}

func SingleHash(in, out chan interface{}) {
	var wg sync.WaitGroup

	for i := range in {
		data := fmt.Sprint(i)
		md5 := DataSignerMd5(data)
		wg.Add(1)

		go func() {
			defer wg.Done()

			ch1 := Crc32Channel(data)
			ch2 := Crc32Channel(md5)

			out <- (<-ch1) + "~" + (<-ch2)
		}()
	}

	wg.Wait()
}

func MultiHash(in, out chan interface{}) {
	var wg sync.WaitGroup

	for i := range in {
		data := fmt.Sprint(i)
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
		result = append(result, fmt.Sprint(i))
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
