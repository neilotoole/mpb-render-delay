package main

import (
	"fmt"
	"time"

	"sync"

	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
)

func main() {
	const delay = time.Second * 3
	start := time.Now()
	delayCh := make(chan struct{})
	time.AfterFunc(delay, func() { close(delayCh) })

	barWg := &sync.WaitGroup{}

	p := mpb.New(mpb.WithAutoRefresh(), mpb.WithRenderDelay(delayCh))

	for i := 0; i < 10; i++ {
		startBarWithLifespan(fmt.Sprintf("%d: aborts before render delay", i), barWg, p, time.Second, true)
	}

	startBarWithLifespan("aborts after render delay", barWg, p, time.Second*6, false)

	barWg.Wait()
	p.Wait()

	fmt.Println("Elapsed: ", time.Since(start))
}

func startBarWithLifespan(name string, wg *sync.WaitGroup, p *mpb.Progress, lifespan time.Duration, colorize bool) {
	death := time.Now().Add(lifespan)
	wg.Add(1)

	if colorize {
		name = "\033[31m" + name + "\033[0m" // red
	}

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
