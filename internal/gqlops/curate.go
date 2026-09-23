package gqlops

import (
	"bytes"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	plainHeadersPattern  = regexp.MustCompile(`(?m)^([ \t]*)plainHeaders[ \t]*$`)
	secretHeadersPattern = regexp.MustCompile(`(?m)^([ \t]*)secretHeaders[ \t]*$`)
)

const devicesAnonymousSelectionSets = 2

func curateCLIContent(path string, content []byte) ([]byte, error) {
	switch filepath.Base(path) {
	case "mutation.accountManagement.disableAccount.txt":
		return replaceRequired(
			content,
			[]byte("disableAccount {"),
			[]byte("disableAccount ( accountId:$accountId ) {"),
			path,
		)
	case "mutation.notification.createSubscriptionGroup.txt",
		"mutation.notification.deleteSubscriptionGroup.txt",
		"mutation.notification.updateSubscriptionGroup.txt",
		"query.notification.txt":
		return addHeaderSelections(content), nil
	case "mutation.xdr.analystFeedback.txt":
		return []byte(`mutation xdrAnalystFeedback($accountId: ID!, $analystFeedbackInput: AnalystFeedbackInput!) {
  xdr(accountId: $accountId) {
    analystFeedback(input: $analystFeedbackInput) {
      story {
        id
      }
    }
  }
}
`), nil
	case "query.devices.txt":
		return removeAnonymousSelectionSets(content)
	default:
		return content, nil
	}
}

func replaceRequired(content, oldValue, newValue []byte, path string) ([]byte, error) {
	if bytes.Count(content, oldValue) != 1 {
		return nil, fmt.Errorf("known repair for %q expected exactly one %q", path, oldValue)
	}
	return bytes.Replace(content, oldValue, newValue, 1), nil
}

func addHeaderSelections(content []byte) []byte {
	updated := plainHeadersPattern.ReplaceAll(content, []byte("${1}plainHeaders {\n${1}\tname\n${1}\tvalue\n${1}}"))
	return secretHeadersPattern.ReplaceAll(updated, []byte("${1}secretHeaders {\n${1}\tname\n${1}}"))
}

func removeAnonymousSelectionSets(content []byte) ([]byte, error) {
	lines := strings.SplitAfter(string(content), "\n")
	output := make([]string, 0, len(lines))
	removed := 0

	for index := 0; index < len(lines); index++ {
		if strings.TrimSpace(lines[index]) != "{" {
			output = append(output, lines[index])
			continue
		}

		depth := 1
		removed++
		for index++; index < len(lines) && depth > 0; index++ {
			depth += strings.Count(lines[index], "{")
			depth -= strings.Count(lines[index], "}")
		}
		index--
		if depth != 0 {
			return nil, fmt.Errorf("unterminated anonymous selection set")
		}
	}

	if removed != devicesAnonymousSelectionSets {
		return nil, fmt.Errorf(
			"devices repair removed %d anonymous selection sets; want %d",
			removed,
			devicesAnonymousSelectionSets,
		)
	}

	return []byte(strings.Join(output, "")), nil
}
