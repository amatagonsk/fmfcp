## copy target

- if `.md` has front matter tag `publish`
- `draft: true` (willing to publish, imo)
  - for hugo user
- not `.md` files


## draft: true && publish (hugo side)
`tag:publish` and not `.md` files are copy.



| publish | ~~draft~~ | copy | memo                              |
| ------- | --------- | ---- | --------------------------------- |
| 1       | 1         | copy | check locally                     |
| 0       | 1         | copy | willing to publish, check locally |
| 1       | 0         | copy |                                   |
| 0       | 0         | skip |                                   |

quartz_v4 not use hugo