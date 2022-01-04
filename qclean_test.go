package qclean_test

import (
	"testing"

	"github.com/po3rin/qclean"
)

var replaceList = map[string]string{
	"ガン":  "がん",
	"前立線": "前立腺",
	"40肩": "四十肩",
	"50肩": "五十肩",
}

func TestCleanSimple(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{
			input: "鳥 が 飛ぶ　動 作",
			want:  "鳥が飛ぶ 動作",
		},
		{
			input: "感染 症　と は",
			want:  "感染症とは",
		},
		{
			input: "誤 嚥 性 肺炎",
			want:  "誤嚥性肺炎",
		},
		{
			input: "がん　を　直す　方法",
			want:  "がんを直す 方法",
		},
		{
			input: "心房 細 動 と は",
			want:  "心房細動とは",
		},
		{
			input: "ｔｒｔ療法",
			want:  "trt療法",
		},
		{
			input: "0 ヶ月 下痢 原因 は",
			want:  "0ヶ月 下痢 原因は",
		},
		{
			input: "1 週間 下痢 と 腹痛",
			want:  "1週間 下痢と腹痛",
		},
		{
			input: "1 歳児 胃腸 炎 症状 下痢 のみ",
			want:  "1歳児 胃腸炎 症状 下痢のみ",
		},
		{
			input: "10ヶ月 しなくなった",
			want:  "10ヶ月 しなくなった",
		},
		{
			input: "11ヶ月 手もみ",
			want:  "11ヶ月 手もみ",
		},
		{
			input: "日本経済新聞 を 読む",
			want:  "日本経済新聞を読む",
		},
		{
			input: "苔 癬 治っ た",
			want:  "苔癬 治った",
		},
	}

	c, err := qclean.NewCleanerWithUserDict("testdata/userdict_test.txt")
	if err != nil {
		t.Fatal(err)
	}

	c.SetReplaceList(replaceList)

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := c.Clean(tt.input)
			if err != nil {
				t.Fatal(err)
			}
			if got != tt.want {
				t.Fatalf("want : %+v\ngot  : %+v", tt.want, got)
			}
		})
	}
}
