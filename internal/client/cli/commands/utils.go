package commands

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"golang.org/x/term"

	"reimagined_eureka/internal/client/infra/logging"
)

// TODO: check on Windows too, and maybe on Linux

func ReadSecretValueMasked(logger logging.ILogger, what string, minSize, maxSize int) (string, error) {
	// lengthHint := getLengthHint(minSize, maxSize)  // TODO: fix "Enter [secretType >=N characters] (%!s(MISSING))"
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

func ReadNonSecretValue(logger logging.ILogger, what string) (string, error) {
	logger.Info("Enter %s: ", what)
	scanner := bufio.NewScanner(os.Stdin)
	var lines []string
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	return strings.Join(lines, "\n"), nil
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

func TrimToNRunes(s string, n int) string {
	runeCount := 0
	for i := range s {
		if runeCount == n {
			return s[:i]
		}
		runeCount++
	}
	return s
}

func ParsePositiveInt(what string) (int, error) {
	rowID, err := strconv.Atoi(what)
	if err != nil || rowID <= 0 {
		return 0, fmt.Errorf("value is a not a positive number")
	}
	return rowID, nil
}
