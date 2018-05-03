package mud

import (
	"bufio"
	"errors"
	"strconv"
	"sync"
	"time"
)

type inputEvent struct {
	inputString string
	err         error
}

const (
	sOUTOFSEQUENCE = iota
	sINESCAPE
	sDIRECTIVE
)

// sleepThenReport is a timeout sequence so that if the escape key is pressed it will register
// as a keypress within a reasonable period of time with the input loop, even if the input
// state machine is in its "inside ESCAPE press listening for extended sequence" state.
func sleepThenReport(stringChannel chan<- inputEvent, myOnce *sync.Once, state *int) {
	time.Sleep(100 * time.Millisecond)

	myOnce.Do(func() {
		*state = sOUTOFSEQUENCE
		stringChannel <- inputEvent{string(rune(27)), nil}
	})
}

func handleKeys(reader *bufio.Reader, stringChannel chan<- inputEvent) {
	inputGone := errors.New("Input ended")
	inEscapeSequence := sOUTOFSEQUENCE
	var myOnce *sync.Once

	for {
		runeRead, _, err := reader.ReadRune()

		if err != nil || runeRead == 3 {
			stringChannel <- inputEvent{"", inputGone}
		}

		if myOnce != nil {
			myOnce.Do(func() { myOnce = nil })
		}

		if inEscapeSequence == sOUTOFSEQUENCE && runeRead == 27 {
			inEscapeSequence = sINESCAPE
			myOnce = new(sync.Once)
			go sleepThenReport(stringChannel, myOnce, &inEscapeSequence)
		} else if inEscapeSequence == sINESCAPE {
			if string(runeRead) == "[" {
				inEscapeSequence = sDIRECTIVE
			} else if runeRead == 27 {
				stringChannel <- inputEvent{string(rune(27)), nil}
			} else {
				inEscapeSequence = sOUTOFSEQUENCE
				if myOnce != nil {
					myOnce.Do(func() { myOnce = nil })
				}
				stringChannel <- inputEvent{string(runeRead), nil}
			}
		} else if inEscapeSequence == sDIRECTIVE {
			switch runeRead {
			case 'A':
				stringChannel <- inputEvent{"UP", nil}
			case 'B':
				stringChannel <- inputEvent{"DOWN", nil}
			case 'C':
				stringChannel <- inputEvent{"RIGHT", nil}
			case 'D':
				stringChannel <- inputEvent{"LEFT", nil}
			default:
				stringChannel <- inputEvent{strconv.QuoteRune(runeRead), nil}
			}
			inEscapeSequence = sOUTOFSEQUENCE
		} else {
			stringChannel <- inputEvent{string(runeRead), nil}
		}
	}
}