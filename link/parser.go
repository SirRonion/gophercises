package htmllinkparser

import (
	"io"
	"strings"

	"golang.org/x/net/html"
)

type Link struct {
	href string
	text string
}

func ParseLinks(f io.Reader) ([]Link, error) {
	links := []Link{}
	link := Link{}
	text := ""
	z := html.NewTokenizer(f)
	depth := 0

	for {
		tt := z.Next()
		switch tt {
		case html.ErrorToken:
			if z.Err() == io.EOF {
				return links, z.Err()
			} else {
				return nil, z.Err()
			}
		case html.StartTagToken:
			tn, hasAttr := z.TagName()
			if string(tn) == "a" {
				depth++
				if depth == 1 {
					if hasAttr {
						key, val, _ := z.TagAttr()
						if string(key) == "href" {
							text := strings.TrimSpace(string(val))
							if len(text) == 0 {
								continue
							}
							link.href = text
						}
					}
				}
			}
		case html.TextToken:
			if depth == 1 {
				text += z.Token().Data + ""
			}
		case html.EndTagToken:
			tn, _ := z.TagName()
			if string(tn) == "a" {
				if depth == 1 {
					depth--
					text = strings.TrimSpace(text)
					link.text = text
					links = append(links, link)
					link = Link{}
					text = ""
				} else {
					depth--
					continue
				}
			}
		case html.CommentToken:
			continue
		}
	}
}
