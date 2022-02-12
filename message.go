package main

type ErrorMessage struct {
	Code    string `json:"error"`
	Message string `json:"message"`
}

var (
	errQueryIsNot5Letters = ErrorMessage{
		Code   : "E001",
		Message: "指定されたクエリが5文字ではありません",
	}

	errInvalidQuery = ErrorMessage{
		Code   : "E002",
		Message: "無効なリクエストです。以下の内容に注意して指定してください。\\n・指定されたリクエストに不正な文字しか存在しなかった\\n・[q] クエリに1文字以上のワイルドカード以外の文字を指定するか、[in] クエリに1文字以上の文字を指定してください",
	}
)