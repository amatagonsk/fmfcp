package main

import (
	"encoding/json"
	"flag"
	"fmt"
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
	draftGj := gjson.Parse(frontMatter).Get("draft")
	var isDraft bool
	err := json.Unmarshal([]byte(draftGj.String()), &isDraft)
	if err != nil {
		isDraft = true
		log.Fatal(err)
	} else {
		isDraft = draftGj.Bool()
	}

	// empty front matter
	if frontMatter == "" {
		return false
	} else {
		tagsGj := gjson.Parse(frontMatter).Get("tags")
		var tags []string
		var isPublishTagContain bool
		err := json.Unmarshal([]byte(tagsGj.String()), &tags)
		if err != nil {
			isPublishTagContain = false
			log.Fatal(err)
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
	// fmt.Println("--------")
	// fmt.Println(filePath)
	return !isPublishCheck(fm), nil
}

func main() {
	src, dest := argCheck()

	//println(src, "\n", dest)
	err := cp.Copy(
		//"C:\\Users\\E14\\Downloads\\tempdir\\tempdir",
		//"C:\\Users\\E14\\Downloads\\tempdir\\piyo",
		src, dest,
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

func argCheck() (string, string) {
	var helpFlag bool
	var cpSrc string
	var cpDest string

	flag.BoolVar(&helpFlag, "h", false, "show help message")
	flag.StringVar(&cpSrc, "src", "", "cp src")
	flag.StringVar(&cpDest, "dest", "", "cp dest")
	flag.Parse()
	if helpFlag {
		printHelp()
	}
	args := flag.Args()
	if len(args) != 2 {
		printHelp()
		log.Fatal("args count wrong")
	}
	cpSrc, cpDest = args[0], args[1]
	return cpSrc, cpDest
}

func printHelp() {
	helpStr := `fmfcp is Front Matter Filter CoPy tool.
check ".md" file & frontmatter contains "tag: publish" or "draft: false" are not copy.
not ".md" file are just copy.
usage: fmfcp $src $dist`
	fmt.Println(helpStr)
	fmt.Println("---")
	flag.PrintDefaults()
}
