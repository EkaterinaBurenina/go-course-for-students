package tagcloud

import "sort"

// TagCloud aggregates statistics about used tags
type TagCloud struct {
	// TODO: add fields if necessary
	//tags []string
	tags_stat map[string]int
}

// TagStat represents statistics regarding single tag
type TagStat struct {
	Tag             string
	OccurrenceCount int
}

// New should create a valid TagCloud instance
// TODO: You decide whether this function should return a pointer or a value
func New() TagCloud {
	// TODO: Implement this
	//stat := TagStat{}
	//tg := TagCloud{tagStat: []stat}
	return TagCloud{tags_stat: make(map[string]int)}
}

// AddTag should add a tag to the cloud if it wasn't present and increase tag occurrence count
// thread-safety is not needed
// TODO: You decide whether receiver should be a pointer or a value
func (tc TagCloud) AddTag(tag string) {
	if _, exist := tc.tags_stat[tag]; exist {
		tc.tags_stat[tag] += 1
	} else {
		tc.tags_stat[tag] = 1
	}
}

// TopN should return top N most frequent tags ordered in descending order by occurrence count
// if there are multiple tags with the same occurrence count then the order is defined by implementation
// if n is greater that TagCloud size then all elements should be returned
// thread-safety is not needed
// there are no restrictions on time complexity
// TODO: You decide whether receiver should be a pointer or a value
func (tc TagCloud) TopN(n int) []TagStat {
	keys := make([]string, 0, len(tc.tags_stat))
	for key := range tc.tags_stat {
		keys = append(keys, key)
	}

	//fmt.Println(basket)
	//fmt.Println("keys before sort", keys)

	sort.SliceStable(keys, func(i, j int) bool {
		return tc.tags_stat[keys[i]] > tc.tags_stat[keys[j]]
	})

	//fmt.Println("keys after sort", keys)

	//res := make([]TagStat, n, n)
	res := []TagStat{}
	for _, k := range keys {
		if len(res) >= n {
			break
		}
		res = append(res, TagStat{Tag: k, OccurrenceCount: tc.tags_stat[k]})
		//fmt.Println("res", res)
		//fmt.Println("len res", len(res))
	}

	return res
}
