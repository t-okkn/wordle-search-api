package main

import (
	"math/rand"
	"net/http"
	"strings"
	"time"

	"wordle-search-api/models"

	"github.com/gin-gonic/gin"
)

// <summary>: 待ち受けるサーバのルーターを定義します
// <remark>: httpHandlerを受け取る関数にそのまま渡せる
func SetupRouter() *gin.Engine {
	router := gin.Default()
	v1 := router.Group("v1")

	v1.GET("/en/search", getEnSearch)
	v1.GET("/en/hint", getEnHint)
	v1.GET("/ja/search", getJaSearch)
	v1.GET("/ja/hint", getJaHint)

	return router
}

// <summary>: 対象のクエリ文字列から絞り込み検索します（Wordle 用）
func getEnSearch(c *gin.Context) {
	EN.getResponses(c, true)
}

// <summary>: Wordle 用のヒントを取得します
func getEnHint(c *gin.Context) {
	EN.getResponses(c, false)
}

// <summary>: 対象のクエリ文字列から絞り込み検索します（WORDLEja 用）
func getJaSearch(c *gin.Context) {
	JA.getResponses(c, true)
}

// <summary>: WORDLWja 用のヒントを取得します
func getJaHint(c *gin.Context) {
	JA.getResponses(c, false)
}

// <summary>: 対象のクエリ文字列から絞り込み検索します
func (v Lang) getResponses(c *gin.Context, isSearch bool) {
	query_prm := c.DefaultQuery("q", "-----")
	include_prm := c.DefaultQuery("in", "")
	exclude_prm := c.DefaultQuery("exclude", "")
	not1_prm := c.DefaultQuery("not1", "")
	not2_prm := c.DefaultQuery("not2", "")
	not3_prm := c.DefaultQuery("not3", "")
	not4_prm := c.DefaultQuery("not4", "")
	not5_prm := c.DefaultQuery("not5", "")
	answer := c.DefaultQuery("answer", "")

	query := v.sanitize(query_prm)
	include := strings.Replace(v.sanitize(include_prm), "-", "", -1)
	exclude := strings.Replace(v.sanitize(exclude_prm), "-", "", -1)
	not1 := strings.Replace(v.sanitize(not1_prm), "-", "", -1)
	not2 := strings.Replace(v.sanitize(not2_prm), "-", "", -1)
	not3 := strings.Replace(v.sanitize(not3_prm), "-", "", -1)
	not4 := strings.Replace(v.sanitize(not4_prm), "-", "", -1)
	not5 := strings.Replace(v.sanitize(not5_prm), "-", "", -1)

	if len([]rune(query)) != 5 || len([]rune(include)) > 5 {
		c.JSON(http.StatusBadRequest, errQueryIsNot5Letters)
		c.Abort()
		return
	}

	if isSearch {
		if query == "-----" && include == "" {
			c.JSON(http.StatusBadRequest, errInvalidQuery)
			c.Abort()
			return
		}
	}

	req := models.RequestSet{
		Query:   query,
		Exclude: exclude,
		Not: map[int]string{
			1: not1,
			2: not2,
			3: not3,
			4: not4,
			5: not5,
		},
	}

	ss := models.SearcherSet{
		Include: include,
	}

	ss.Searcher = make(map[int][]rune, 5)
	s := v.getSearcher(req)

	for key, val := range s {
		ss.Searcher[key] = make([]rune, 0, len(s[key]))

		for _, r := range val {
			if r != DELETED {
				ss.Searcher[key] = append(ss.Searcher[key], r)
			}
		}
	}

	result := v.getResult(ss, answer)

	if !isSearch {
		src := rand.NewSource(time.Now().UnixNano())
		rnd := rand.New(src)

		hint := result[rnd.Intn(len(result))]
		result = []string{hint}
	}

	res := models.ResponseSet{
		Query:   query,
		Include: separateString(include),
		Result:  result,
		Options: map[string]string{
			"exclude": separateString(exclude),
			"not1":    separateString(not1),
			"not2":    separateString(not2),
			"not3":    separateString(not3),
			"not4":    separateString(not4),
			"not5":    separateString(not5),
		},
	}

	c.JSON(http.StatusOK, res)
}

func (v Lang) sanitize(input string) string {
	var sb strings.Builder
	sb.Grow(len(input))

	list := []rune{}
	if v == EN {
		list = AllowEN
	} else if v == JA {
		list = AllowJA
	}

	for _, c := range input {
		if contains(c, list) {
			sb.WriteRune(c)
		} else {
			sb.WriteString("-")
		}
	}

	return sb.String()
}

func (v Lang) getSearcher(req models.RequestSet) map[int][]rune {
	s := make(map[int][]rune, 5)
	done := []bool{false, false, false, false, false}

	// ----- 準備フェーズ -----
	list := []rune{}
	if v == EN {
		list = AllowEN
	} else if v == JA {
		list = AllowJA
	}

	s[1] = make([]rune, len(list))
	s[2] = make([]rune, len(list))
	s[3] = make([]rune, len(list))
	s[4] = make([]rune, len(list))
	s[5] = make([]rune, len(list))

	for key := range s {
		for pos, i := range list {
			s[key][pos] = i
		}
	}

	count_deleted := func(r []rune) int {
		count := 0

		for _, i := range r {
			if i == DELETED {
				count += 1
			}
		}

		return count
	}

	// ----- クエリについての処理フェーズ -----
	if req.Query == "-----" {
		for pos, qchar := range req.Query {
			if qchar == HYPHEN {
				continue
			}

			for cnt, rval := range s[pos+1] {
				if qchar != rval {
					s[pos+1][cnt] = DELETED
				}
			}

			done[pos] = true
		}
	}

	// ----- 文字単位の除外についての処理フェーズ -----
	for key, not := range req.Not {
		if not == "" {
			continue
		}
		if done[key-1] {
			continue
		}

		for cnt, rval := range s[key] {
			for _, i := range not {
				if rval == i {
					s[key][cnt] = DELETED
				}
			}
		}

		d := count_deleted(s[key])
		if d == 0 || d == 1 {
			done[key-1] = true
		}
	}

	// ----- excludeについての処理フェーズ -----
	if req.Exclude != "" {
		for key := range s {
			if done[key-1] {
				continue
			}

			for cnt, rval := range s[key] {
				for _, i := range req.Exclude {
					if rval == i {
						s[key][cnt] = DELETED
					}
				}
			}
		}
	}

	return s
}

func (v Lang) getResult(ss models.SearcherSet, answer string) []string {
	ans := strings.ToLower(answer)
	include := []rune(ss.Include)

	list := []string{}
	if v == EN {
		if ans == "1" || ans == "true" || ans == "yes" || ans == "on" {
			list = EWordsAns
		} else {
			list = EWordsAll
		}

	} else if v == JA {
		list = JWords
	}

	want := make([]string, 0, len(list))

	for _, word := range list {
		for pos, char := range word {
			if !contains(char, include) {
				continue
			}

			if contains(char, ss.Searcher[pos+1]) {
				want = append(want, word)
			}
		}
	}

	return want
}

func contains(r rune, list []rune) bool {
	if len(list) == 0 {
		return false
	}

	for _, v := range list {
		if r == v {
			return true
		}
	}

	return false
}

func separateString(input string) string {
	if input == "" { return "" }

	words_count := len([]rune(input))
	if words_count == 1 {
		return input
	}

	var sb strings.Builder
	sb.Grow(len(input) + words_count)

	for pos, i := range input {
		if pos != 0 {
			sb.WriteString(", ")
		}

		sb.WriteRune(i)
	}

	return sb.String()
}
