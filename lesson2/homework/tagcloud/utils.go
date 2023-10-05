package tagcloud

func FindStringInTagCloud(cloud *TagCloud, value string) *TagStat {

	for i := range cloud.stats {
		if cloud.stats[i].Tag == value {
			return &cloud.stats[i]
		}
	}
	return nil
}
