package utils

import "strings"

type CaseChanger interface {
	Uppercase(str string) string
	Lowercase(str string) string
}

type SpaceTrimmer interface {
	TrimSpaces(str string) string
}

type MyTextConverter struct {
}

func (c *MyTextConverter) Uppercase(str string) string {
	return strings.ToUpper(str)
}

func (c *MyTextConverter) Lowercase(str string) string {
	return strings.ToLower(str)
}

func (c *MyTextConverter) TrimSpaces(str string) string {
	return strings.TrimSpace(str)
}
