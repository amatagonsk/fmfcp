package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"path/filepath"

	cp "github.com/otiai10/copy"

	_ "flag"
	"log"
	"os"

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

/* `tag: publish` or `draft: false` */
func isPublishCheck(frontMatter string) bool {
	// empty front matter
	if frontMatter == "" {
		return false
	} else {
		draftGj := gjson.Parse(frontMatter).Get("draft")
		var isDraft bool
		err := json.Unmarshal([]byte(draftGj.String()), &isDraft)
		if err != nil {
			isDraft = false
			// log.Panic(err)
		} else {
			isDraft = draftGj.Bool()
		}

		tagsGj := gjson.Parse(frontMatter).Get("tags")
		var tags []string
		var isPublishTagContain bool
		err = json.Unmarshal([]byte(tagsGj.String()), &tags)
		if err != nil {
			isPublishTagContain = false
			log.Panic(err)
		} else {
			isPublishTagContain = slices.Contains(tags, "publish")
		}

		isPublish := !(isDraft || isPublishTagContain)
		if isPublish {
			return false
		} else {
			return true
		}
	}
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
	if err := md.Convert([]byte(bytes), io.Discard, parser.WithContext(ctx)); err != nil {
		return false, fmt.Errorf("convert error: %s\n%w", filePath, err)
	}

	var fm string
	if data := frontmatter.Get(ctx); data != nil {
		var meta map[string]any
		if err := data.Decode(&meta); err != nil {
			return false, fmt.Errorf("decode error: %s\n%w", filePath, err)
		}

		formatted, err := json.MarshalIndent(meta, "", "  ")
		if err != nil {
			return false, fmt.Errorf("format error: %s\n%w", filePath, err)
		}

		fm = string(formatted)
	}
	return !isPublishCheck(fm), nil
}

func main() {
	src, dest := argCheck()

	//println(src, "\n", dest)
	err := cp.Copy(
		src, dest,
		cp.Options{
			Skip: func(info os.FileInfo, src, dest string) (bool, error) {
				if info.IsDir() || filepath.Ext(src) != ".md" {
					fmt.Println(src)
					return false, nil
				} else {
					isSkip, err := fileFilter(filepath.Join(src))
					if !isSkip {
						fmt.Println(src)
					}
					return isSkip, err
				}
			},
		},
	)
	if err != nil {
		log.Panic(err)
	}
}

func argCheck() (string, string) {
	var helpFlag bool
	var cpSrc string
	var cpDest string

	flag.BoolVar(&helpFlag, "h", false, "show help message")
	flag.BoolVar(&helpFlag, "help", false, "show help message")
	flag.StringVar(&cpSrc, "src", "", "cp source")
	flag.StringVar(&cpDest, "dest", "", "cp destination")
	flag.Parse()
	if helpFlag {
		printHelpAndExit()
	}
	args := flag.Args()
	if len(args) != 2 {
		printHelpAndExit()
	}
	cpSrc, cpDest = args[0], args[1]
	return cpSrc, cpDest
}

func printHelpAndExit() {
	helpStr := `fmfcp is Front Matter Filter CoPy tool.
check ".md" file & frontmatter contains "tag: publish" or "draft: true" are copy.
not ".md" file are just copy.

usage: fmfcp $src $dest`
	fmt.Println(helpStr)
	// fmt.Println("---")
	// flag.PrintDefaults()
	os.Exit(0)
}
