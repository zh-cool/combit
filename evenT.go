package main

import (
	"github.com/ethereum/go-ethereum/event"
	"fmt"
	"time"
)

type MyEvent struct {
	X int
	Y int
}

func wevent(ev *event.TypeMux) {
	fmt.Println("entry wevent")
	sub := ev.Subscribe(MyEvent{})

	sch := sub.Chan()
	fmt.Println("Range sch")
	for obj := range sch {
		switch eve := obj.Data.(type){
		case MyEvent:
			fmt.Println("Recv MyEvent\n", eve)

		}
	}
	fmt.Println("Out wevent")
}

func nowait(ev *event.TypeMux) {
	fmt.Println("entry nowait")
	ev.Subscribe(MyEvent{})

	time.Sleep(1*time.Second)
	fmt.Println("ouat nowait")
}

func evfeed(ev *event.Feed) {
	ch := make(chan MyEvent)
	sub := ev.Subscribe(ch)

	for {
		select {
		case event := <-ch:
			fmt.Println(event)

			// Err() channel will be closed when unsubscribing.
		case <-sub.Err():
			return
		}
	}

}

func init() {
	fmt.Println("Austin")
}

func main(){
	ev := new(event.Feed)
	var ve event.Feed

	go evfeed(ev)
	go evfeed(ev)

	go evfeed(&ve)
	go evfeed(&ve)

	time.Sleep(1*time.Second)

	ev.Send(MyEvent{1,1})
	ve.Send(MyEvent{2,2})

	time.Sleep(1*time.Second)
}
