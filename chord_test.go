package chord

import (
	"testing"
	"time"
)

var (
	responses            ResponsesMapType
	sendMessagesFromCode BackgroundProcessType
)

func TestCreateAndStartBot(t *testing.T) {
	// responds to Hello, HELLO, hello, hElLo
	responses = make(ResponsesMapType)
	responses["hello"] = func() string { return "hi, I'm a bot" }

	// entry point for background process with a channel to send messages to discord
	sendMessagesFromCode = func(messages chan<- string) {
		for {
			time.Sleep(10 * time.Second)
			messages <- "Something happened in your code"
		}
	}

	// Run your bot!!
	_, err := NewChord(
		"", // token (no need to add "Bot" prefix, we'll do that for you)
		"", // guild
		"", // channel (usually "general")
		responses,
		sendMessagesFromCode,
	)
	if err != nil {
		t.Error(err)
	}
}
