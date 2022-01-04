package qclean_test

import (
	"strings"
	"testing"

	"github.com/po3rin/qclean"
)

var replaceList = map[string]string{
	"ガン":   "がん",
	"前立線":  "前立腺",
	"４０肩":  "四十肩",
	"５０肩":  "五十肩",
	"ＰC R": "PCR",
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
			input: "日本経済新聞 を 読む",
			want:  "日本経済新聞を読む",
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

func TestSelectJoinedRaw(t *testing.T) {
	tests := []struct {
		name      string
		raw       []string
		converted []string
		want      []string
	}{
		{
			name:      "simple1",
			raw:       []string{"大豆", "製品", "取りすぎ"},
			converted: []string{"大豆", "製品", "取り", "すぎ"},
			want:      []string{"大豆", "製品", "取りすぎ"},
		},
		{
			name:      "simple2",
			raw:       []string{"大豆", "取りすぎ"},
			converted: []string{"大豆", "取り", "すぎ", "原因"},
			want:      []string{"大豆", "取りすぎ", "原因"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := qclean.SelectJoinedRaw(tt.raw, tt.converted)
			if strings.Join(got, " ") != strings.Join(tt.want, " ") {
				t.Fatalf("want : %+v\ngot  : %+v", tt.want, got)
			}
		})
	}
}
