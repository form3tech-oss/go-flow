package option

import "github.com/form3tech-oss/go-flow/pkg/types"

type optionState struct {
	channelBuffer int
}
type Option func(state *optionState)

func BufferedChannel(capacity int) Option {
	return func(state *optionState) {
		state.channelBuffer = capacity
	}
}

func getOptions(options ...Option) *optionState {
	state := &optionState{}
	for _, option := range options {
		option(state)
	}
	return state
}

func CreateChannel(options ...Option) chan types.Element {
	state := getOptions(options...)
	if state.channelBuffer > 0 {
		return make(chan types.Element, state.channelBuffer)
	}
	return make(chan types.Element)
}
