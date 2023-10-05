package main

import "lecture02_homework/tagcloud"

func main() {

	cloud := tagcloud.New()

	cloud.AddTag("1")
	cloud.AddTag("1")

	for _, tag := range cloud.TopN(10) {
		print("Tag: ", tag.Tag, " Count: ", tag.OccurrenceCount)
	}
}
