package main

import (
	"github.com/gin-gonic/gin"
)

// summary => 待ち受けるサーバのルーターを定義します
// return::*gin.Engine =>
// remark => httpHandlerを受け取る関数にそのまま渡せる
/////////////////////////////////////////
func SetupRouter() *gin.Engine {
	router := gin.Default()
	v1 := router.Group("v1")

	v1.GET("/search/:query", getSearch)
	v1.GET("/hint", getHint)

	return router
}

// summary => 対象のクエリ文字列から絞り込み検索します
// param::c => [p] gin.Context構造体
/////////////////////////////////////////
func getSearch(c *gin.Context) {
	//途中で処理を中断
	//c.Abort()
}

// summary => ヒントを取得します
// param::c => [p] gin.Context構造体
/////////////////////////////////////////
func getHint(c *gin.Context) {
	//途中で処理を中断
	//c.Abort()
}