package newt

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

const (
	TuiStandardMenuHelp = `Enter menu letter and id`
	TuiStandardMenu     = `Menu [a]dd, [m]odify, [r]emove, [q]uit (making changes)`
)

// getDigit get a numeric answer for a string input that is
// greater than or equal than zero and less than the length of the list.
// NOTE: This returns a zero based array position.
func getDigit(buf *bufio.Reader, list []string) (int, bool) {
    answer := getAnswer(buf, "", true)
	if answer != "" {
		pos, err := strconv.Atoi(answer)
		if err == nil {
			// Adjust input to zero based array address.
			pos--
			if pos >= 0 && pos < len(list) {
				return pos, true
			}
		}
	}
	return -1, false
}

// getAnswer get a Y/N response from buffer
func getAnswer(buf *bufio.Reader, defaultAnswer string, lower bool) string {
	answer, err := buf.ReadString('\n')
	if err != nil {
		return ""
	}
	answer = strings.TrimSpace(answer)
	if answer == "" {
		answer = defaultAnswer
	}
	if lower {
		return strings.ToLower(answer)
	}
	return answer
}

// getAnswers returns an answer which has an initial verb and an predicate separted
// by a space. E.g. "modify 1" -> "modify" "1"
func getAnswers(buf *bufio.Reader, defaultAnswer string, defaultValue string, lower bool) (string, string) {
	var (
		answer1 string
		answer2 string
	)
	rawAnswer := getAnswer(buf, defaultAnswer, false)
	if strings.Contains(rawAnswer, " ") {
		parts := strings.SplitN(rawAnswer, " ", 2)
		answer1, answer2 = parts[0], parts[1]
	} else {
		answer1 = rawAnswer
	}
	answer1 = strings.TrimSpace(answer1)
	answer2 = strings.TrimSpace(answer2)
	if answer1 == "" {
		answer1 = defaultAnswer
	}
	if lower {
		return strings.ToLower(answer1), answer2
	}
	return answer1, answer2
}

// selectMenuItem displays a description, a list of menu items (selected by name or number)
// returns the selected menu action and optional modify using getAnswers().
func selectMenuItem(in io.Reader, out io.Writer, topMsg string, bottomMsg string, list []string, numberedList bool, defaultAnswer string, defaultValue string, lower bool) (string, string) {
	readBuffer := bufio.NewReader(in)
	fmt.Fprintf(out, "%s\n\n", topMsg)
	for i, item := range list {
		if numberedList {
			fmt.Fprintf(out, "\t%d: %s\n", (i + 1), item)
		} else {
			fmt.Fprintf(out, "\t%s\n", item)
		}
	}
	fmt.Fprintf(out, "\n%s\n", bottomMsg)
	return getAnswers(readBuffer, defaultAnswer, defaultValue, lower)
}
