// Elasticsearch DSL
// @Author: Golion
// @Date: 2017.3

package elastic

import (
	"encoding/json"
	"fmt"
	"strings"

	"vidmate.com/common/utils"
)

type esPutShards struct {
	Succ int `json:"successful"`
}

type esPutResult struct {
	Shards esPutShards `json:"_shards"`
}

func (putResult *esPutResult) Init() {
	putResult.Shards.Succ = 0
}

func (putResult *esPutResult) CheckPutSucc() bool {
	if putResult.Shards.Succ >= 1 {
		return true
	} else {
		return false
	}
}

type Term struct {
	Field string
	Value string
}

func GetBaseScrollQuery(size int, terms []Term, sources []string) (string, error) {
	var termQueries []string
	for i := 0; i < len(terms); i++ {
		termQueries = append(termQueries, `{ "term" : { "`+terms[i].Field+`" : "`+terms[i].Value+`" } }`)
	}
	js, err := utils.JSONDecode(string(`{
		"size": ` + fmt.Sprintf("%v", size) + `,
		"query": {
			"bool": {
				"must": [ ` + strings.Join(termQueries, ",") + ` ]
			}
		},
		"_source": ["` + strings.Join(sources, `", "`) + `"]
	}`))
	if err != nil {
		return "", fmt.Errorf("[GetBaseScrollQuery] error=[%v]", err.Error())
	}
	output, err := json.Marshal(js)
	if err != nil {
		return "", fmt.Errorf("[GetBaseScrollQuery] error=[%v]", err.Error())
	}
	return string(output), nil
}

type ESBaseHit struct {
	Index string  `json:"_index"`
	Type  string  `json:"_type"`
	Id    string  `json:"_id"`
	Score float64 `json:"_score"`
}

type ESBaseHits struct {
	Total    int         `json:"total"`
	MaxScore float64     `json:"max_score"`
	Hits     []ESBaseHit `json:"hits"`
}

type ESScrollResult struct {
	ScrollId string     `json:"_scroll_id"`
	Hits     ESBaseHits `json:"hits"`
}

type ESScrollQuery struct {
	ScrollTime string `json:"scroll"`
	ScrollId   string `json:"scroll_id"`
}
