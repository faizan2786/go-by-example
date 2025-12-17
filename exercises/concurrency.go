package main

import (
	"fmt"
	"sync"
	"time"
)

type token struct{} // empty struct used for signalling (not to carry any data)

type SafeGroup struct {
	wg        sync.WaitGroup
	semaphore chan token // a channel to control concurrency limit (i.e. a semaphore)
}

func (sg *SafeGroup) Wait() { sg.wg.Wait() }

func (sg *SafeGroup) SetLimit(n int) {
	if n <= 0 {
		panic("limit must be > 0")
	}
	sg.semaphore = make(chan token, n) // initialise the limit channel with n number of buffer
}

func (sg *SafeGroup) Go(work func()) {

	// check if limit is set
	if sg.semaphore != nil {
		// add a token before starting the go routine
		// hence, when the limit channel is full, the send will block
		// and no further go routines will be spawned until the channel can
		// start receiving again (i.e. a token is remvoed by a finished go routine)
		sg.semaphore <- token{}
	}

	sg.wg.Add(1)
	go func() {
		defer sg.wg.Done()
		work()
		if sg.semaphore != nil {
			<-sg.semaphore // remove a token from the channel when finished
		}
	}()
}

func main() {
	var sg SafeGroup
	sg.SetLimit(10)
	for i := range 40 {
		sg.Go(func() {
			time.Sleep(time.Second)
			fmt.Printf("%d.", i)
		})
	}
	sg.Wait()
	fmt.Println("\nmain done.")
}

func main_non_limit_test() {

	var sg SafeGroup
	for i := range 3 {
		sg.Go(func() {
			time.Sleep(time.Duration(i) * time.Second)
			fmt.Printf("worker %d done\n", i)
		})
	}
	sg.Wait()
	fmt.Println("main done.")
}
