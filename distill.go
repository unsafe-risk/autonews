package autonews

import (
	"golang.org/x/net/html"
)

var usefulattrs = map[string]bool{
	"src":            true,
	"href":           true,
	"alt":            true,
	"title":          true,
	"role":           true,
	"aria-label":     true,
	"aria-hidden":    true,
	"aria-atomic":    true,
	"focusable":      true,
	"class":          true,
	"id":             true,
	"name":           true,
	"type":           true,
	"value":          true,
	"placeholder":    true,
	"checked":        true,
	"disabled":       true,
	"readonly":       true,
	"selected":       true,
	"required":       true,
	"for":            true,
	"tabindex":       true,
	"maxlength":      true,
	"minlength":      true,
	"pattern":        true,
	"size":           true,
	"min":            true,
	"max":            true,
	"step":           true,
	"multiple":       true,
	"autocomplete":   true,
	"autofocus":      true,
	"form":           true,
	"formaction":     true,
	"formenctype":    true,
	"formmethod":     true,
	"formtarget":     true,
	"formnovalidate": true,
}

func distillPipeline(n *html.Node) {
	if n == nil {
		return
	}

	unusefulE := false
	for i := range n.Attr {
		if !usefulattrs[n.Attr[i].Key] {
			unusefulE = true
			break
		}
	}

	if unusefulE {
		newAttr := make([]html.Attribute, 0, len(n.Attr))
		for i := range n.Attr {
			if usefulattrs[n.Attr[i].Key] {
				newAttr = append(newAttr, n.Attr[i])
			}
		}
		n.Attr = newAttr
	}

L:
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.ElementNode {
			switch c.Data {
			case "svg":
				c.FirstChild = nil
				c.LastChild = nil
				continue L
			case "script", "style", "link", "noscript":
				defer n.RemoveChild(c)
				continue L
			}
		}
		distillPipeline(c)
	}
}
