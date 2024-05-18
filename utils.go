package autonews

import "github.com/lemon-mint/vermittlungsstelle/llm"

func getTexts(c *llm.Content) []string {
	var texts []string
	if c == nil {
		return texts
	}

	for _, t := range c.Parts {
		switch t.Type() {
		case llm.SegmentTypeText:
			texts = append(texts, string(t.(llm.Text)))
		}
	}

	return texts
}
