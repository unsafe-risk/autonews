package autonews

import (
	"context"
	"fmt"
	"strings"

	"github.com/lemon-mint/vermittlungsstelle/llm"
)

const newstitle_prompt = `Read this newsletter and come up with a different title based on the topic you think your readers will be most interested in!:

<news>
%s
</news>

After reading, identify the main features and points covered in the document.

Then write a friendly and interesting one-sentence title for the article that addresses those key points.

Make your sentences as vivid as possible, and you can even use exclamation marks.

The title should use soft, positive language to convey the essence of the document in an easily understandable way.

Avoid using string formats like Markdown and HTML.

Write the title within the <title> tag. The title must be written in English.`

func generateNewsTitle(ctx context.Context, model llm.LLM, document string) (string, error) {
	response := model.GenerateStream(ctx, &llm.ChatContext{}, &llm.Content{
		Role: llm.RoleModel,
		Parts: []llm.Segment{
			llm.Text(fmt.Sprintf(newstitle_prompt, document)),
		},
	})
	for range response.Stream {
	} // Wait for the stream to finish.

	if response.Err != nil {
		return "", response.Err
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
