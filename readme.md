# Front Matter Filter CoPy

obsidian is personal knowledge bases, contains secret & publishable.
front matter `tag:publish` files're copy, others're not. filter copy tool.
so not split vault & able to use other markdown service. (hugo, quartz, etc.)

---

fmfcp is Front Matter Filter CoPy tool.
check ".md" file & frontmatter contains "tag: publish" or "draft: true" are copy.
not ".md" file are just copy.


## install

`go install github.com/amatagonsk/go-fmfcp@latest`


## usage(help)

`fmfcp $src $dest`
