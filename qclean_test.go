package qclean_test

import (
	"testing"

	"github.com/po3rin/qclean"
)

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
			input: "苔 癬 治っ た",
			want:  "苔癬 治った",
		},
		{
			input: "がん　を　直す　方法",
			want:  "がんを直す 方法",
		},
		{
			input: "心房 細 動 と は",
			want:  "心房細動とは",
		},
	}

	c, err := qclean.NewCleaner()
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
