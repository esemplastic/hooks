package hooks

import (
	"fmt"
)

var msgHub = NewHub()

type message struct {
	from string
	to   string
	body string
}

func newMessage(from, to, body string) message {
	fmt.Println("newMessage")
	return message{
		from: from,
		to:   to,
		body: body,
	}
}

// sender part

func send(from, to, body string) {
	msgHub.RunFunc(send, from, to, body)
}

// receiver part

func filterMessage(msg message) error {
	fmt.Printf("filtering message: %#v\n", msg)
	if msg.from == "from name" {
		return nil
	}
	return fmt.Errorf("bad message")
}

func printAnyErr(err error) {
	if err != nil {
		fmt.Println(err.Error())
	}
}

func sendValidMessage(msg message) {
	fmt.Println("Sending message...")
	fmt.Printf("%#v\n", msg)
}

func ExampleChain() {

	// flow :
	// -> newMessage(from,to,body) -> filterMessage(message) -> receive(message)
	// |< if Any error then the printAnyErr will be executed and the chain will be stopped.

	msgHub.RegisterFunc(send, Chain(newMessage, filterMessage, sendValidMessage, printAnyErr))
	send("from name", "to name", "message contents")

	// Output:
	// newMessage
	// filtering message: hooks.message{from:"from name", to:"to name", body:"message contents"}
	// Sending message...
	// hooks.message{from:"from name", to:"to name", body:"message contents"}
}
