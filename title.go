package autonews

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/lemon-mint/coord/llm"
)

const title_prompt = `Please carefully read the following document:

<document>
%s
</document>

After reading, identify the main features and points covered in the document.

Then write a friendly and interesting one-sentence title for the article that addresses those key points.

Make your sentences as vivid as possible, and you can even use exclamation marks.

The title should use soft, positive language to convey the essence of the document in an easily understandable way.

Avoid using string formats like Markdown and HTML.

Write the title within the <title> tag. The title must be written in English.`

var ErrTitleFailed = errors.New("title failed")

func generateTitle(ctx context.Context, model llm.LLM, document string) (string, error) {
	response := model.GenerateStream(
		ctx,
		&llm.ChatContext{},
		llm.TextContent(llm.RoleUser, fmt.Sprintf(title_prompt, document)),
	)

	err := response.Wait()
	if err != nil {
		return "", err
	}

	texts := getTexts(response.Content)
	if len(texts) == 0 {
		return "", ErrTitleFailed
	}

	text := strings.Join(texts, "")

	_, after, ok := strings.Cut(text, "<title>")
	if !ok {
		return "", ErrTitleFailed
	}

	title, _, ok := strings.Cut(after, "</title>")
	if !ok {
		return "", ErrTitleFailed
	}

	return strings.TrimSpace(title), nil
}
