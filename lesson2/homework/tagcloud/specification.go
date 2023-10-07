package tagcloud

import (
	"sort"
)

// TagCloud aggregates statistics about used tags
type TagCloud struct {
	stats []TagStat
}

// TagStat represents statistics regarding single tag
type TagStat struct {
	Tag             string
	OccurrenceCount int
}

// New should create a valid TagCloud instance
func New() TagCloud {
	return TagCloud{}
}

// AddTag should add a tag to the cloud if it wasn't present and increase tag occurrence count
// thread-safety is not needed
func (cloud *TagCloud) AddTag(tag string) {
	found := FindStringInTagCloud(cloud, tag)

	if found != nil {
		found.OccurrenceCount += 1
	} else {
		cloud.stats = append(cloud.stats, TagStat{Tag: tag, OccurrenceCount: 1})
	}
}

// TopN should return top N most frequent tags ordered in descending order by occurrence count
// if there are multiple tags with the same occurrence count then the order is defined by implementation
// if n is greater that TagCloud size then all elements should be returned
// thread-safety is not needed
// there are no restrictions on time complexity
func (cloud *TagCloud) TopN(n int) []TagStat {
	sort.SliceStable(cloud.stats, func(i, j int) bool { return cloud.stats[i].OccurrenceCount > cloud.stats[j].OccurrenceCount })

	if n >= len(cloud.stats) {
		return cloud.stats
	} else {
		return cloud.stats[:n]
	}
}
