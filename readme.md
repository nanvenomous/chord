# Description

A declarative way to spin up a simple discord bot with a map of possible responses and a single background process for sending unprompted messages.

# Installation
> go get github.com/mrgarelli/chord

### Dependencies
> go get github.com/bwmarrin/discordgo

# Usage
[see example in test file](https://github.com/mrgarelli/chord/blob/master/chord_test.go)

add information for your bot
```
	_, err := NewChord(
		"", // token (no need to add "Bot" prefix, we'll do that for you)
		"", // guild
		"", // channel (usually "general")
```

> go mod tidy

> go test
