package autonews

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/lemon-mint/vermittlungsstelle/llm"
)

const translation_prompt = `Translate the document into Korean.

<document>
%s
</document>

Always use a friendly and bright tone.
Do not distort the meaning of the sentence.
Do not translate names or abbreviations.
Do not translate yaml keys.

Write the translated content inside the <korean> block.`

var ErrTranslationFailed = errors.New("translation failed")

func generateTranslation(ctx context.Context, model llm.LLM, document string) (string, error) {
	response := model.GenerateStream(ctx, &llm.ChatContext{}, &llm.Content{
		Role: llm.RoleModel,
		Parts: []llm.Segment{
			llm.Text(fmt.Sprintf(translation_prompt, document)),
		},
	})
	for range response.Stream {
	} // Wait for the stream to finish.

	if response.Err != nil {
		return "", response.Err
	}

	texts := getTexts(response.Content)
	if len(texts) == 0 {
		return "", ErrTranslationFailed
	}

	text := strings.Join(texts, "")

	_, after, ok := strings.Cut(text, "<korean>")
	if !ok {
		return "", ErrTranslationFailed
	}

	translation, _, ok := strings.Cut(after, "</korean>")
	if !ok {
		return "", ErrTranslationFailed
	}

	return strings.TrimSpace(translation), nil
}
