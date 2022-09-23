package security

import (
	"github.com/microcosm-cc/bluemonday"
	"regexp"
)

var SafeHTMLPolicy *bluemonday.Policy

func init() {
	SafeHTMLPolicy = bluemonday.UGCPolicy()
	SafeHTMLPolicy.AllowStyles("color", "background-color", "text-align", "width", "height", "font-size", "font-weight", "padding-left").Globally()
	SafeHTMLPolicy.AddTargetBlankToFullyQualifiedLinks(true)
	SafeHTMLPolicy.AllowElements("iframe")
	SafeHTMLPolicy.AllowAttrs("width").Matching(bluemonday.Number).OnElements("iframe")
	SafeHTMLPolicy.AllowAttrs("height").Matching(bluemonday.Number).OnElements("iframe")
	SafeHTMLPolicy.AllowAttrs("src").OnElements("iframe")
	SafeHTMLPolicy.AllowAttrs("frameborder").Matching(bluemonday.Number).OnElements("iframe")
	SafeHTMLPolicy.AllowAttrs("allow").Matching(regexp.MustCompile(`[a-z; -]*`)).OnElements("iframe")
	SafeHTMLPolicy.AllowAttrs("allowfullscreen").OnElements("iframe")
}
