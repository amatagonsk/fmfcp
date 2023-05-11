package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	_ "github.com/go-playground/validator/v10"
	"github.com/tidwall/gjson"
	_ "github.com/tidwall/gjson"
	"github.com/yuin/goldmark"
	_ "github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	_ "github.com/yuin/goldmark/parser"
	"go.abhg.dev/goldmark/frontmatter"
	_ "go.abhg.dev/goldmark/frontmatter"
	"golang.org/x/exp/slices"
	_ "golang.org/x/exp/slices"
)

const TES = `---
creation date: 2022-06-17 16:55
modification date: Thursday 30th June 2022 16:16:55
# tags:
# - publish
# - jq
# draft: false
hi: greeting
---

## head2
new news is here

tag: publish
draft: false
`

/* `tag: publish` or `draft: false` */
func publishFilter(meta string) bool {
	//println(meta)

	draftGj := gjson.Get(meta, "draft")
	tagsGj := gjson.Get(meta, "tags")

	var isDraft bool
	err := json.Unmarshal([]byte(draftGj.String()), &isDraft)
	if err != nil {
		isDraft = true
		err = nil
	}
	//println("isDraft: ", isDraft)

	var tags []string
	var isTagContain bool
	err = json.Unmarshal([]byte(tagsGj.String()), &tags)
	if err != nil {
		isTagContain = false
	} else {
		isTagContain = slices.Contains(tags, "publish")
	}

	for i, s := range tags {
		fmt.Println("tag", i, s)
	}
	//err = slices.Contains(tags, "publish")

	isPublish := !isDraft || isTagContain
	return isPublish
}

func main() {

	//res, err := run(&req)
	res, err := run(TES)
	if err != nil {
		res = &response{Error: err.Error()}
	}
	//_ = res

	println(publishFilter(res.Frontmatter))

}

type response struct {
	HTML        string
	Frontmatter string
	Error       string
}

func run(str string) (*response, error) {
	var formats []frontmatter.Format

	md := goldmark.New(
		goldmark.WithExtensions(
			&frontmatter.Extender{
				Formats: formats,
			},
		),
	)

	ctx := parser.NewContext()
	var buf bytes.Buffer
	if err := md.Convert([]byte(str), &buf, parser.WithContext(ctx)); err != nil {
		return nil, fmt.Errorf("convert markdown: %w", err)
	}

	var fm string
	if data := frontmatter.Get(ctx); data != nil {
		var meta map[string]any
		if err := data.Decode(&meta); err != nil {
			return nil, fmt.Errorf("decode frontmatter: %w", err)
		}

		formatted, err := json.MarshalIndent(meta, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("format frontmatter: %w", err)
		}

		fm = string(formatted)
	}

	return &response{
		HTML:        buf.String(),
		Frontmatter: fm,
	}, nil
}
