package main

import (
	"sort"
	"strconv"
	"strings"
)

func ExecutePipeline(jobs ...job) {

}

func SingleHash(in, out chan interface{}) {
	outData := <-out
	in <- DataSignerCrc32(outData.(string)) + "~" + DataSignerCrc32(DataSignerMd5(outData.(string)))
}

func MultiHash(in, out chan interface{}) {
	outData := <-out
	for i := 0; i <= 5; i++ {
		in <- DataSignerCrc32(strconv.Itoa(i) + outData.(string))
	}
}

func CombineResults(in, out chan interface{}) {
	var result []string

	for i := range out {
		result = append(result, i.(string))
	}

	sort.Strings(result)
	in <- strings.Join(result, "_")
}
