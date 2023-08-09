package main

import (
	"html"
	"regexp"
)

var (
	detailRegexp = regexp.MustCompile(`<a href="(.+)".+target=`)
	nameRegexp   = regexp.MustCompile(`<div\s+class="name".+?>(.+?)</div>`)
	// 品质
	rarityRegexp = regexp.MustCompile(`品质.+>(.+)</span>`)

	onsealRegexp = regexp.MustCompile(`sell_num_show="(\d+)"`)
)

func parseList(h string) (links []string) {
	result := detailRegexp.FindAllStringSubmatch(h, -1)
	for _, v := range result {
		if len(v) != 2 {
			continue
		}
		links = append(links, html.UnescapeString(v[1]))
	}
	return
}
