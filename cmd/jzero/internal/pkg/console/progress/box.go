package progress

import (
	"fmt"

	"github.com/jzero-io/jzero/cmd/jzero/internal/pkg/console"
)

// StageState tracks the console state for a progress stage.
type StageState struct {
	shown        bool
	hasErrorItem bool
}

// ConsumeStage renders progress messages until done.
func ConsumeStage(progressChan <-chan Message, done <-chan struct{}, title string, quiet, headerShown bool) StageState {
	state := StageState{
		shown: headerShown,
	}

	for {
		select {
		case msg, ok := <-progressChan:
			if !ok {
				return state
			}
			renderMessage(msg, title, quiet, &state)
		case <-done:
			for {
				select {
				case msg, ok := <-progressChan:
					if !ok {
						return state
					}
					renderMessage(msg, title, quiet, &state)
				default:
					return state
				}
			}
		}
	}
}

// FinishStage renders the footer and any error detail lines.
func FinishStage(title string, quiet bool, state *StageState, err error) {
	if quiet {
		return
	}

	if err != nil {
		if !state.hasErrorItem {
			if item := ItemFromError(err); item != "" {
				if !state.shown {
					fmt.Printf("%s\n", console.BoxHeader("", title))
					state.shown = true
				}
				fmt.Printf("%s\n", console.BoxErrorItem(item))
				state.hasErrorItem = true
			}
		}

		if !state.shown {
			fmt.Printf("%s\n", console.BoxHeader("", title))
			state.shown = true
		}

		for _, line := range console.NormalizeErrorLines(err.Error()) {
			fmt.Printf("%s\n", console.BoxDetailItem(line))
		}
	}

	if !state.shown {
		return
	}

	if err != nil {
		fmt.Printf("%s\n\n", console.BoxErrorFooter())
		return
	}

	fmt.Printf("%s\n\n", console.BoxSuccessFooter())
}

func renderMessage(msg Message, title string, quiet bool, state *StageState) {
	if quiet {
		return
	}

	if !state.shown {
		fmt.Printf("%s\n", console.BoxHeader("", title))
		state.shown = true
	}

	switch msg.Type {
	case TypeFile:
		fmt.Printf("%s\n", console.BoxItem(msg.Value))
	case TypeError:
		fmt.Printf("%s\n", console.BoxErrorItem(msg.Value))
		state.hasErrorItem = true
	case TypeDebug:
		// Debug command lines are intentionally not rendered in progress boxes.
	}
}
