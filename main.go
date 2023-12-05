package main

import (
	"fmt"
	"time"

	"github.com/vbauerster/mpb/v8"
)

func main() {
	const delay = time.Second * 5
	start := time.Now()
	delayCh := make(chan struct{})
	time.AfterFunc(delay, func() { close(delayCh) })

	p := mpb.New(mpb.WithAutoRefresh(), mpb.WithRenderDelay(delayCh))
	bar := p.New(0, mpb.BarStyle(), mpb.BarRemoveOnComplete())
	bar.Abort(true)
	bar.Wait() // <-- blocks here until render delay completes
	p.Wait()

	fmt.Println("Elapsed: ", time.Since(start))
}
