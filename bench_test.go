package qclean_test

import (
	"testing"

	"github.com/po3rin/qclean"
)

func BenchmarkClean(b *testing.B) {
	c, err := qclean.NewCleanerWithUserDict("testdata/userdict_test.txt")
	if err != nil {
		b.Fatal(err)
	}
	c.SetReplaceList(replaceList)

	b.ResetTimer()
	// Nはコマンド引数から与えられたベンチマーク時間から自動で計算される
	for i := 0; i < b.N; i++ {
		got, _ := c.Clean("抗 ガン 剤 やめる と どうなる")
		_ = got
	}
}
