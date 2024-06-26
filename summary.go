package autonews

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/lemon-mint/coord/llm"
)

const summary_prompt = `Please carefully read the following document:

<document>
%s
</document>

After reading, identify the key features and main points covered in the document.

Then, write a friendly 3-paragraph summary of the document that covers those key points.

Make your sentences as lively as possible, you can even use exclamation points.

Don't start with "This document is".

Always start your summary by explaining what the project or document does.

Include the insightful points of the document.

The summary should use soft, positive language to convey the essence of the document in an easily understandable way.

Never use string formatting such as Markdown and HTML.

If there is no appropriate content to summarize or the given document is an error page, print "NO_CONTENT"

Write your summary inside <summary> tags. The summary should be in English.`

var ErrSummaryNotProvided = errors.New("summary not provided")

func generateSummary(ctx context.Context, model llm.LLM, document string) (string, error) {
	response := model.GenerateStream(
		ctx,
		&llm.ChatContext{},
		llm.TextContent(llm.RoleUser, fmt.Sprintf(summary_prompt, document)),
	)

	err := response.Wait()
	if err != nil {
		return "", err
	}

	texts := getTexts(response.Content)
	if len(texts) == 0 {
		return "", ErrSummaryNotProvided
	}

	text := strings.Join(texts, "")
	if strings.Contains(text, "NO_CONTENT") {
		return "", ErrSummaryNotProvided
	}

	_, after, ok := strings.Cut(text, "<summary>")
	if !ok {
		return "", ErrSummaryNotProvided
	}

	summary, _, ok := strings.Cut(after, "</summary>")
	if !ok {
		return "", ErrSummaryNotProvided
	}

	summary = strings.TrimSpace(summary)
	summary = strings.ReplaceAll(summary, "\n\n", "\n")
	summary = strings.ReplaceAll(summary, "\n\n", "\n")
	summary = strings.ReplaceAll(summary, "\n", "\n\n")

	lines := strings.Split(summary, "\n")
	var LinesWithoutEmpty []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) > 0 {
			LinesWithoutEmpty = append(LinesWithoutEmpty, line)
		}
	}
	summary = strings.Join(LinesWithoutEmpty, "\n\n")
	summary = strings.TrimSpace(summary)

	return summary, nil
}
