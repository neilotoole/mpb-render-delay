package main

import (
	"fmt"
	"time"

	"sync"

	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
)

func main() {
	const delay = time.Second * 5
	start := time.Now()
	delayCh := make(chan struct{})
	time.AfterFunc(delay, func() { close(delayCh) })

	barWg := &sync.WaitGroup{}

	p := mpb.New(mpb.WithAutoRefresh(), mpb.WithRenderDelay(delayCh))

	for i := 0; i < 10; i++ {
		startBarWithLifespan(fmt.Sprintf("%d: aborts before render delay", i), barWg, p, time.Second)
	}

	startBarWithLifespan("aborts after render delay", barWg, p, time.Second*10)

	barWg.Wait()
	p.Wait()

	fmt.Println("Elapsed: ", time.Since(start))
}

func startBarWithLifespan(name string, wg *sync.WaitGroup, p *mpb.Progress, lifespan time.Duration) {
	death := time.Now().Add(lifespan)
	wg.Add(1)

	bar := p.New(0,
		mpb.BarStyle(),
		mpb.BarWidth(40),
		mpb.PrependDecorators(decor.Name(name, decor.WCSyncWidthR)),
		mpb.AppendDecorators(decor.CurrentNoUnit("%d")),
		mpb.BarRemoveOnComplete(),
	)
	go func() {
		defer wg.Done()
		for {

			if time.Now().After(death) {
				bar.Abort(true)
				return
			}
			bar.IncrBy(1)
			time.Sleep(time.Millisecond * 100)
		}
	}()
}
