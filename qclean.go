package qclean

import (
	"strings"

	ipa "github.com/ikawaha/kagome-dict-ipa-neologd"
	"github.com/ikawaha/kagome-dict/dict"

	// "github.com/ikawaha/kagome-dict/ipa"
	"github.com/ikawaha/kagome/v2/tokenizer"
	"github.com/y-bash/go-gaga"
)

var Dict = ipa.Dict()

type Cleaner struct {
	tknz        *tokenizer.Tokenizer
	norm        *gaga.Normalizer
	replaceList map[string]string
}

// NewCleaner ユーザー辞書を使わずに初期化する
func NewCleaner() (*Cleaner, error) {
	tknz, err := tokenizer.New(Dict, tokenizer.OmitBosEos())
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

func (c *Cleaner) ApplyReplace(txt string) string {
	for k, v := range c.replaceList {
		txt = strings.ReplaceAll(txt, k, v)
	}
	return txt
}

func (c *Cleaner) Clean(txt string) (string, error) {
	txt = strings.ReplaceAll(txt, "　", " ")
	txt = strings.ReplaceAll(txt, "\n", "")

	txt = c.ApplyReplace(txt)
	txt = c.Norm(txt)

	rawSplit := strings.Split(txt, " ")
	if len(rawSplit) == 0 {
		return txt, nil
	}

	if len(rawSplit) <= 2 {
		return txt, nil
	}

	txt = strings.ReplaceAll(txt, " ", "")
	txt = c.Norm(txt)

	tokens := c.tknz.Tokenize(txt)

	var prefix_pool string
	var next_join bool
	var results []string
	for _, t := range tokens {
		var pos string
		if len(t.Features()) >= 6 {
			pos = strings.Join(t.Features()[:7], ",")
		} else {
			pos = strings.Join(t.Features()[:len(t.Features())], ",")
		}

		if len(results) > 0 &&
			(strings.Contains(pos, "副詞,助詞類接続") ||
				strings.Contains(pos, "助詞,連体化") ||
				strings.Contains(pos, "助詞,格助詞") ||
				strings.Contains(pos, "助詞,接続助詞") ||
				strings.Contains(pos, "助詞,副助詞") ||
				strings.Contains(pos, "助詞,並立助詞") ||
				strings.Contains(pos, "助詞,係助詞") ||
				strings.Contains(pos, "助詞,副詞化") ||
				strings.Contains(pos, "名詞,非自立,一般,*,*,*,の") ||
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
			(strings.Contains(pos, "名詞,接尾,助数詞") ||
				// strings.Contains(pos, "動詞,自立,*,*,五段・ラ行,連用タ接続") ||
				strings.Contains(pos, "名詞,接尾,特殊,*,*,*") ||
				strings.Contains(pos, "動詞,自立,*,*,五段・ラ行,基本形") ||
				strings.Contains(pos, "形容詞,非自立") ||
				strings.Contains(pos, "助動詞,*,*,*,特殊") ||
				strings.Contains(pos, "動詞,自立,*,*,五段・ラ行,連用形") ||

				// Special case ...
				// There may be a problem in this case.
				strings.Contains(pos, "名詞,接尾,一般,*,*,*,児") ||
				strings.Contains(pos, "カスタム接尾")) {
			results[len(results)-1] = results[len(results)-1] + t.Surface
			continue
		}

		if strings.Contains(pos, "接頭詞") || t.Surface == "-" || strings.Contains(pos, "形容詞,非自立,*,*,形容詞・アウオ段,連用テ接続") {
			prefix_pool = t.Surface
			continue
		}

		results = append(results, prefix_pool+t.Surface)

		// reset
		prefix_pool = ""
	}

	results = SelectJoinedRaw(rawSplit, results)
	return strings.Join(results, " "), nil
}

func (c *Cleaner) CleanAll(txts []string) ([]string, error) {
	results := make([]string, 0, 0)
	for _, t := range txts {
		result, err := c.Clean(t)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	return results, nil
}

// SelectJoinedRaw 元のクエリでトークン分解されてしまっているものがある場合、元の形を採用する
func SelectJoinedRaw(raw []string, converted []string) []string {
	rawMap := make(map[string]struct{}, len(raw))
	for _, r := range raw {
		rawMap[r] = struct{}{}
	}

	for i := 0; i < len(converted)-1; i++ {
		compare := converted[i] + converted[i+1]
		_, ok := rawMap[compare]
		if !ok {
			continue
		}
		converted = append(converted[:i], converted[i+1:]...)
		converted[i] = compare
		i++
	}
	return converted
}
