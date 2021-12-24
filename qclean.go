package qclean

import (
	"strings"

	ipa "github.com/ikawaha/kagome-dict-ipa-neologd"
	"github.com/ikawaha/kagome-dict/dict"

	// "github.com/ikawaha/kagome-dict/ipa"
	"github.com/ikawaha/kagome/v2/tokenizer"
	"github.com/y-bash/go-gaga"
)

type Cleaner struct {
	tknz        *tokenizer.Tokenizer
	norm        *gaga.Normalizer
	replaceList map[string]string
}

func NewCleaner() (*Cleaner, error) {
	tknz, err := tokenizer.New(ipa.Dict(), tokenizer.OmitBosEos())
	if err != nil {
		return nil, err
	}

	norm, err := gaga.Norm(gaga.LatinToNarrow | gaga.KanaToWide)
	if err != nil {
		return nil, err
	}

	return &Cleaner{
		tknz: tknz,
		norm: norm,
	}, nil
}

func (c *Cleaner) SetReplaceList(replaceList map[string]string) {
	c.replaceList = replaceList
}

// NewCleanerWithUserDict ユーザー辞書ファイルとともに初期化する
func NewCleanerWithUserDict(filepath string) (*Cleaner, error) {
	udic, err := dict.NewUserDict(filepath)
	if err != nil {
		return nil, err
	}

	tknz, err := tokenizer.New(ipa.Dict(), tokenizer.UserDict(udic), tokenizer.OmitBosEos())
	if err != nil {
		return nil, err
	}

	norm, err := gaga.Norm(gaga.LatinToNarrow | gaga.KanaToWide)
	if err != nil {
		return nil, err
	}

	return &Cleaner{
		tknz: tknz,
		norm: norm,
	}, nil
}

func (c *Cleaner) Norm(txt string) string {
	txt = c.norm.String(txt)
	// txt = string(norm.NFKC.Bytes([]byte(txt)))
	return txt
}

func (c *Cleaner) Clean(txt string) (string, error) {
	rawSplit := strings.Split(txt, " ")
	if len(rawSplit) == 0 {
		return txt, nil
	}

	txt = strings.ReplaceAll(txt, " ", "")
	txt = strings.ReplaceAll(txt, "　", "")

	for k, v := range c.replaceList {
		txt = strings.ReplaceAll(txt, k, v)
	}

	txt = c.Norm(txt)

	tokens := c.tknz.Tokenize(txt)

	var prefix_pool string
	var next_join bool
	var results []string

	for _, t := range tokens {
		var pos string
		if len(t.Features()) >= 6 {
			pos = strings.Join(t.Features()[:6], ",")
		} else {
			pos = strings.Join(t.Features()[:len(t.Features())], ",")
		}

		if len(results) > 0 &&
			(strings.Contains(pos, "副詞,助詞類接続") ||
				strings.Contains(pos, "助詞,連体化") ||
				strings.Contains(pos, "助詞,格助詞") ||
				strings.Contains(pos, "助詞,接続助詞") ||
				strings.Contains(pos, "動詞,自立,*,*,サ変・スル,未然形")) {
			results[len(results)-1] = results[len(results)-1] + t.Surface
			next_join = true
			continue
		}
		if len(results) > 0 && next_join {
			results[len(results)-1] = results[len(results)-1] + t.Surface
			next_join = false
			continue
		}

		if len(results) > 0 &&
			(strings.Contains(pos, "名詞,接尾") ||
				strings.Contains(pos, "詞,自立,*,*,五段・ラ行,基本形") ||
				strings.Contains(pos, "助動詞,*,*,*,特殊")) {
			results[len(results)-1] = results[len(results)-1] + t.Surface
			continue
		}

		if strings.Contains(pos, "接頭詞") || t.Surface == "-" {
			prefix_pool = t.Surface
			continue
		}

		results = append(results, prefix_pool+t.Surface)
		prefix_pool = ""
	}

	results = SelectJoinedRaw(rawSplit, results)
	return strings.Join(results, " "), nil
}

// SelectJoinedRaw 元々のクエリで分解されてしまっているものは元の形を採用する
func SelectJoinedRaw(raw []string, converted []string) []string {
	result := make([]string, 0)
	addedmap := make(map[string]struct{})

	for _, c := range converted {
		var checkcnt int
		for _, r := range raw {
			if strings.Contains(r, c) {
				if _, ok := addedmap[r]; ok {
					break
				}
				addedmap[r] = struct{}{}
				result = append(result, r)
				break
			}
			checkcnt++
		}
		if checkcnt == len(raw) {
			if _, ok := addedmap[c]; ok {
				continue
			}
			addedmap[c] = struct{}{}
			result = append(result, c)
		}

	}
	return result
}
