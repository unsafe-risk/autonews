package autonews

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/lemon-mint/vermittlungsstelle/llm"
)

const summary_prompt = `Please carefully read the following document:

<document>
%s
</document>

After reading, identify the key features and main points covered in the document. 

Then, write a friendly 2-paragraph summary of the document that covers those key points.
The summary should use soft, positive language to convey the essence of the document in an easily understandable way.

Never use string formatting such as Markdown and HTML.

If there is no appropriate content to summarize or the given document is an error page, print "NO_CONTENT"

Write your summary inside <summary> tags. The summary should be in English.`

func generateSummaryPrompt(document string) string {
	return fmt.Sprintf(summary_prompt, document)
}

var ErrSummaryNotProvided = errors.New("summary not provided")

func generateSummary(ctx context.Context, model llm.LLM, document string) (string, error) {
	response := model.GenerateStream(ctx, &llm.ChatContext{}, &llm.Content{
		Role: llm.RoleModel,
		Parts: []llm.Segment{
			llm.Text(generateSummaryPrompt(document)),
		},
	})
	for range response.Stream {
	} // Wait for the stream to finish.

	if response.Err != nil {
		return "", response.Err
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

	return strings.TrimSpace(summary), nil
}
