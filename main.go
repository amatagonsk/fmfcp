package main

import (
	"encoding/json"
	"fmt"

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
	"log"
	"os"
	"path/filepath"
)

/* `tag: publish` or `draft: false` */
func isPublishCheck(frontMatter string) bool {
	//println(frontMatter)
	draftGj := gjson.Parse(frontMatter).Get("draft")
	var isDraft bool
	err := json.Unmarshal([]byte(draftGj.String()), &isDraft)
	if err != nil {
		isDraft = true
		log.Fatal(err)
	} else {
		isDraft = draftGj.Bool()
	}

	tagsGj := gjson.Parse(frontMatter).Get("tags")
	var tags []string
	var isTagContain bool
	err = json.Unmarshal([]byte(tagsGj.String()), &tags)
	if err != nil {
		isTagContain = false
		log.Fatal(err)
	} else {
		isTagContain = slices.Contains(tags, "publish")
	}
	isPublish := !isDraft || isTagContain
	return isPublish
}

func fileFilter(filePath string) (bool, error) {
	bytes, err := os.ReadFile(filePath)
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
	// isPublishCheck: true:publish, false:draft
	// copySkip: true:skip, false:copy
	return !isPublishCheck(fm), nil
}

func main() {
	err := cp.Copy(
		"C:\\Users\\E14\\Downloads\\tempdir\\tempdir",
		"C:\\Users\\E14\\Downloads\\tempdir\\piyo",
		cp.Options{
			Skip: func(info os.FileInfo, src, dest string) (bool, error) {
				if info.IsDir() || filepath.Ext(src) != ".md" {
					return false, nil
				} else {
					isSkip, err := fileFilter(filepath.Join(src))
					return isSkip, err
				}
			},
		},
	)
	if err != nil {
		log.Fatal(err)
	}
}
