package htmllinkparser

import (
	"fmt"
	"io"
	"os"
	"testing"
)

func TestParser(t *testing.T) {
	cases := []struct {
		file     string
		expected []Link
	}{
		{
			file: "ex1.html",
			expected: []Link{
				{
					href: "/other-page",
					text: "A link to another page",
				},
			},
		},
		{
			file: "ex2.html",
			expected: []Link{
				{
					href: "https://www.twitter.com/joncalhoun",
					text: "Check me out on twitter",
				},
				{
					href: "https://github.com/gophercises",
					text: "Gophercises is on Github",
				},
			},
		},
		{
			file: "ex3.html",
			expected: []Link{
				{
					href: "#",
					text: "Login",
				},
				{
					href: "/lost",
					text: "Lost? Need help?",
				},
				{
					href: "https://twitter.com/marcusolsson",
					text: "@marcusolsson",
				},
			},
		},
		{
			file: "ex4.html",
			expected: []Link{
				{
					href: "/dog-cat",
					text: "dog cat",
				},
			},
		},
	}

	for _, c := range cases {
		f, err := os.Open(c.file)
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()
		links, err := ParseLinks(f)
		if err != io.EOF {
			t.Fatal(err)
		}

		for i := range c.expected {
			if links[i] != c.expected[i] {
				t.Errorf("FAIL:\nexpected: %v\ngot: %v\n", c.expected[i], links[i])
			} else {
				fmt.Printf("PASS:\nexpected: %v\ngot: %v\n", c.expected[i], links[i])
			}
		}
	}
}
