package main

import (
	"bus-practice/pkg/bus"
	"flag"
	"fmt"
	"log"
	"time"
)

func main() {
	busMemSize := flag.Int("m", 10, "amount of bus last messages stored in memeory")
	flag.Parse()

	b := bus.NewBus(*busMemSize)
	r1, err := b.Connect(bus.ConnectOpts{
		Stream: "out",
		Type:   "r1",
	})
	go func(r *bus.Reader) {
		for m := range r.ReadCh() {
			log.Println(fmt.Sprintf("reader 1 msg: %v", m))
		}
	}(r1.(*bus.Reader))
	if err != nil {
		log.Println("reader 1 connect fail")
	}

	f1, err := b.Connect(bus.ConnectOpts{
		Stream: "in",
		Type:   "r1",
		Name:   "feeder 1",
	})

	f1.(*bus.Feeder).Feed() <- &bus.FeedMsg{
		Type: "r1",
		Msg:  "nice f1",
	}
	f1.(*bus.Feeder).Feed() <- &bus.FeedMsg{
		Type: "r1",
		Msg:  "nice 2 f1",
	}

	f1.(*bus.Feeder).Feed() <- &bus.FeedMsg{
		Type: "r2",
		Msg:  "nice 2 f1",
	}

	time.Sleep(time.Second * 3)
}
