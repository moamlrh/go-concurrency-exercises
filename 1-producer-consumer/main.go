//////////////////////////////////////////////////////////////////////
//
// Given is a producer-consumer scenario, where a producer reads in
// tweets from a mockstream and a consumer is processing the
// data. Your task is to change the code so that the producer as well
// as the consumer can run concurrently
//

package main

import (
	"fmt"
	"sync"
	"time"
)

func producer(stream Stream, c chan<- *Tweet) {
	fmt.Println("Producer started")
	defer close(c)
	for {
		tweet, err := stream.Next()
		if err == ErrEOF {
			return
		}
		c <- tweet
	}
}

func consumer(c <-chan *Tweet, wg *sync.WaitGroup) {
	defer wg.Done()
	for t := range c {
		if t.IsTalkingAboutGo() {
			fmt.Println(t.Username, "\ttweets about golang with id ")
		} else {
			fmt.Println(t.Username, "\tdoes not tweet about golang with id ")
		}
	}
}

func main() {
	start := time.Now()
	stream := GetMockStream()
	c := make(chan *Tweet)
	var wg sync.WaitGroup

	go producer(stream, c)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go consumer(c, &wg)
	}

	wg.Wait()

	fmt.Printf("Process took %s\n", time.Since(start))
}
