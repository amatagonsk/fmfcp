package main

import (
	_ "bytes"
	"encoding/json"
	"fmt"
	_ "github.com/go-playground/validator/v10"
	cp "github.com/otiai10/copy"
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
	"os"
)

const TES = `---
creation date: 2022-06-17 16:55
modification date: Thursday 30th June 2022 16:16:55
# tags:
# - publish
# - jq
draft: true
hi: greeting
---

## head2
new news is here

end
`

/* `tag: publish` or `draft: false` */
func publishFilter(front string) bool {

	//println(front)
	//println("---")

	draftGj := gjson.Parse(front).Get("draft")
	var isDraft bool

	err := json.Unmarshal([]byte(draftGj.String()), &isDraft)

	if err != nil {
		println("no draft")
		isDraft = true
	} else {
		//draftGj.String()
		println(" draft!!!")
		//isDraft = true
		isDraft = draftGj.Bool()
	}

	tagsGj := gjson.Parse(front).Get("tags")
	var tags []string
	var isTagContain bool
	err = json.Unmarshal([]byte(tagsGj.String()), &tags)
	if err != nil {
		println("!!!!!tags not found!!!!!!!")
		isTagContain = slices.Contains(tags, "publish")
	} else {
		println("????????????")
		isTagContain = slices.Contains(tags, "publish")
	}

	for i, s := range tags {
		fmt.Println("tags:", i, s)
	}

	//isPublish := !isDraft
	isPublish := !isDraft || isTagContain

	//println("isDraft:", isDraft)
	//println("isTagContain:", isTagContain)
	//println("isPublish:", isPublish)
	return isPublish
}

func fileFilter(fileInfo string) (bool, error) {
	//println(fileInfo)
	bytes, err := os.ReadFile(fileInfo)
	if err != nil {
		return false, fmt.Errorf("err open file: %w", err)
	}

	var formats []frontmatter.Format

	md := goldmark.New(
		goldmark.WithExtensions(
			&frontmatter.Extender{
				Formats: formats,
			},
		),
	)

	ctx := parser.NewContext()
	//var bb io.Writer
	if err := md.Convert([]byte(bytes), os.Stdout, parser.WithContext(ctx)); err != nil {
		//if err := md.Convert([]byte(bytes), os.Stdout); err != nil {
		return false, fmt.Errorf("convert markdown: %w", err)
	}

	var fm string
	if data := frontmatter.Get(ctx); data != nil {
		var meta map[string]any
		if err := data.Decode(&meta); err != nil {
			return false, fmt.Errorf("decode frontmatter: %w", err)
		}

		formatted, err := json.MarshalIndent(meta, "", "  ")
		if err != nil {
			return false, fmt.Errorf("format frontmatter: %w", err)
		}

		fm = string(formatted)
	}
	// true:skip, false:copy
	return !publishFilter(fm), nil
}

func main() {
	////res, err := run(&req)
	//res := publishFilter(TES)
	//
	//println(res)
	////if err != nil {
	////	res = &response{Error: err.Error()}
	////}
	////_ = res

	err := cp.Copy(
		"C:\\Users\\E14\\Downloads\\tempdir\\tempdir",
		"C:\\Users\\E14\\Downloads\\tempdir\\piyo",
		cp.Options{
			Skip: func(info os.FileInfo, src, dest string) (bool, error) {
				isSkip, err := fileFilter(src)
				return isSkip, err
			},
		},
	)
	fmt.Println(err) // nil

}
