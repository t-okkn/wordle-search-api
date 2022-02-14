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
		Message: "無効なリクエストです。指定されたリクエストに不正な文字しか存在しなかった可能性があります。もしくは、[q] クエリに1文字以上のワイルドカード以外の文字を指定するか、[in] クエリに1文字以上の文字を指定してください。",
	}

	errConflictOptions = ErrorMessage{
		Code   : "E101",
		Message: "指定したオプションどうしが競合しているため、検索ができません。指定したオプションを見直してください。",
	}
)