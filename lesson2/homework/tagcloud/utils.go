package tagcloud

func FindStringInTagCloud(cloud *TagCloud, value string) *TagStat {

	for _, stat := range cloud.stats {
		if stat.Tag == value {
			return &stat
		}
	}

	return nil
}
