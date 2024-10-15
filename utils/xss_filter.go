package utils

import "github.com/microcosm-cc/bluemonday"

// XssFilter Xss 过滤器
func XssFilter(str string) string {
	p := bluemonday.NewPolicy()
	p.AllowElements("br", "img")
	p.AllowAttrs("style", "src", "alt", "width", "height", "loading").OnElements("img")
	return p.Sanitize(str)
}
