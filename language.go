package main

import "strings"

type Lang string

const (
	NA Lang = "NA"
	EN Lang = "EN"
	JA Lang = "JA"
)

func (v Lang) String() string {
	return string(v)
}

func ParseLang(s string) (v Lang) {
	switch strings.ToLower(s) {
	case "en", "english", "wordle", "えいご", "エイゴ", "英語":
		v = EN

	case "ja", "japanese", "wordleja", "にほんご", "ニホンゴ", "日本語":
		v = JA

	default:
		v = NA
	}

	return
}
