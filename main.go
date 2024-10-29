package main

import (
	"os"
	"sync"
	"time"

	g "xabbo.b7c.io/goearth"
)

var ext = g.NewExt(g.ExtInfo{
    Title: "G-Terminal-Test",
    Author: "b7",
    Version: "1.0",
})

func main() {
    var once sync.Once
    ext.Activated(func() {
        once.Do(func() {
            ch := make(chan int)
            go send(ch)
            go recv(ch)
        })
    })
    ext.Run()
}

// Writes 1024 bytes to stdout and sends the number of bytes written to the channel every 100ms.
// Reports the first error received from Write to G-Earth's extension console.
func send(ch chan<- int) {
    buf := make([]byte, 1024)
    for i := range buf {
        buf[i] = 'X'
    }

    errored := false
    for range time.Tick(100 * time.Millisecond) {
        n, err := os.Stdout.Write(buf)
        ch <- n
        if err != nil && !errored {
            errored = true
            ext.Log(err)
        }
    }
}

// Receives the number of bytes written from the channel and reports the results to G-Earth's extension console.
// If nothing was received after 1 second, reports that a lockup was detected with the total number of bytes written.
func recv(ch <-chan int) {
    total := 0
    for {
        select {
        case n := <-ch:
            total += n
            ext.Logf("%d (%d)", n, total)
        case <-time.After(time.Second):
            ext.Logf("lockup detected!")
            ext.Logf("total bytes written: %d", total)
            return
        }
    }
}

