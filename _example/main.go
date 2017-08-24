package main

import (
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/mix3/parawri"
)

func main() {
	ps := parawri.NewParallelStdout()
	kv := map[string]io.Writer{
		"w1": ps.NewAppendWriter(),
		"w2": ps.NewAppendWriter(),
		"w3": ps.NewAppendWriter(),
	}

	for _, k := range []string{"w1", "w2", "w3"} {
		fmt.Fprintf(kv[k], "%s ... ", k)
	}

	var wg sync.WaitGroup
	for sec, k := range []string{"w1", "w2", "w3"} {
		wg.Add(1)
		v := kv[k]
		go func(sec int, k string, v io.Writer) {
			defer wg.Done()
			time.Sleep(time.Second * time.Duration(sec+1))
			fmt.Fprint(v, color.GreenString("done"))
		}(sec, k, v)
	}
	wg.Wait()
}
