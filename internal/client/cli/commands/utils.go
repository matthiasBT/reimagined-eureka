package commands

import (
	"bufio"
	"fmt"
	"os"

	"golang.org/x/term"

	"reimagined_eureka/internal/client/infra/logging"
)

func readSecretValueMasked(logger logging.ILogger, what string, minSize, maxSize int) (string, error) {
	// lengthHint := getLengthHint(minSize, maxSize)  // TODO: fix "Enter [what >=N characters] (%!s(MISSING))"
	logger.Info("Enter %s: ", what)
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return "", fmt.Errorf("failed to read %s: %v", what, err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	reader := bufio.NewReader(os.Stdin)
	var input []rune
	for {
		r, _, err := reader.ReadRune()
		if err != nil {
			return "", fmt.Errorf("failed to read %s: %v", what, err)
		}
		switch r {
		case '\r', '\n':
			logger.Warning("\n\r")
			result := string(input)
			if minSize != 0 && len(result) < minSize {
				return "", fmt.Errorf("%s is shorter than %d characters", what, minSize)
			} else if maxSize != 0 && len(result) > maxSize {
				return "", fmt.Errorf("%s is longer than %d characters", what, maxSize) // TODO: use maybe?
			}
			return result, nil
		case '\x7f', '\b': // Backspace key
			if len(input) > 0 {
				logger.Warning("\b \b") // Move back, write space to clear, and move back again
				input = input[:len(input)-1]
			}
		default:
			logger.Warning("*")
			input = append(input, r)
		}
	}
}

func getLengthHint(minSize, maxSize int) string {
	var lengthHint string
	if minSize != 0 && maxSize != 0 {
		lengthHint = fmt.Sprintf("%d-%d characters", minSize, maxSize)
	} else if minSize == 0 && maxSize != 0 {
		lengthHint = fmt.Sprintf("<=%d characters", maxSize)
	} else if minSize != 0 && maxSize == 0 {
		lengthHint = fmt.Sprintf(">=%d characters", minSize)
	} else {
		lengthHint = "any length"
	}
	return lengthHint
}
