package autonews

import (
	"context"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
	"github.com/lemon-mint/vermittlungsstelle/llm"
	"golang.org/x/net/html"
)

func NewAutoNews(model llm.LLM) *AutoNews {
	return &AutoNews{
		model: model,
	}
}

type AutoNews struct {
	model llm.LLM
}

func (g *AutoNews) Crawl(ctx context.Context, url string) (html string, err error) {
	r := rod.New()
	defer r.Close()

	err = r.Connect()
	if err != nil {
		return
	}

	page, err := r.Page(proto.TargetCreateTarget{
		URL: url,
	})
	if err != nil {
		return
	}
	page = page.Context(ctx)
	page = page.Timeout(time.Second * 30)
	err = page.WaitDOMStable(time.Second, 1)
	if err != nil {
		return
	}

	html, err = page.HTML()
	return
}

func (g *AutoNews) GetSummary(ctx context.Context, url string) (summary string, err error) {
	orightml, err := g.Crawl(ctx, url)
	if err != nil {
		return
	}

	htmlnode, err := html.Parse(strings.NewReader(orightml))
	if err != nil {
		return
	}
	distillPipeline(htmlnode)

	var sb strings.Builder
	err = html.Render(&sb, htmlnode)
	if err != nil {
		return
	}

	summary, err = generateSummary(ctx, g.model, sb.String())
	return
}
