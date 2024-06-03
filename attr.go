package sqlcache

import (
	"regexp"
	"strconv"
)

var (
	attrRegexp = regexp.MustCompile(`(@cache-ttl|@cache-max-rows|@cache-query-string) (\d+)`)
)

type attributes struct {
	ttl              int
	maxRows          int
	cacheQueryString bool
}

func getAttrs(query string) *attributes {
	matches := attrRegexp.FindAllStringSubmatch(query, -1)

	var attrs attributes
	for _, match := range matches {
		if len(match) != 3 {
			return nil
		}
		switch match[1] {
		case "@cache-ttl":
			ttl, _ := strconv.Atoi(match[2])
			attrs.ttl = ttl
		case "@cache-max-rows":
			maxRows, _ := strconv.Atoi(match[2])
			attrs.maxRows = maxRows
		case "@cache-query-string":
			cacheQueryString, _ := strconv.Atoi(match[2])
			attrs.cacheQueryString = cacheQueryString != 0
		}
	}

	if attrs.ttl == 0 || attrs.maxRows == 0 {
		return nil
	}

	return &attrs
}
