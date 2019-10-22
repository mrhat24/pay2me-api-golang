package pay2me_api

import (
	"crypto/md5"
	"encoding/hex"
	"net/http"
	"sort"
)

const API_KEY_EMPTY = "API_KEY_EMPTY"

const API_KEY_ERROR = "API_KEY_ERROR"

const INVALID_SIGNATURE = "INVALID_SIGNATURE"

type CreateDealParams struct {
	OrderID     string `json:"order_id"`
	OrderDesc   string `json:"order_desc"`
	OrderAmount string `json:"order_amount"`
	Signature   string `json:"signature"`
}

type Deal struct {
	CreateDate  string `json:"create_date" schema:"create_date"`
	OrderID     string `json:"order_id" schema:"order_id"`
	UpdateDate  string `json:"update_date" schema:"update_date"`
	ObjectID    string `json:"object_id" schema:"object_id"`
	Redirect    string `json:"redirect" schema:"redirect"`
	OrderAmount string `json:"order_amount" schema:"order_amount"`
	Signature   string `json:"signature" schema:"signature"`
	ExpireDate  string `json:"expire_date" schema:"expire_date"`
	OrderDesc   string `json:"order_desc" schema:"order_desc"`
	Status      string `json:"status" schema:"status"`
}

type Pay2MeParams map[string]string

func (params Pay2MeParams) Sorted() Pay2MeParams {
	newParams := Pay2MeParams{}
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	s := sort.StringSlice{}
	s = keys
	sort.Strings(s)
	for _, k := range s {
		newParams[k] = (params)[k]
	}
	return newParams
}

func (params Pay2MeParams) InlineValues() string {
	result := ""
	for k := range params {
		result += params[k]
	}
	return result
}

type Pay2MeApi struct {
	key string
}

func (p *Pay2MeApi) doRequest(r *http.Request) (*http.Response, error) {
	c := http.DefaultClient
	return c.Do(r)
}

func (p *Pay2MeApi) getSignature(params Pay2MeParams, key string) string {
	sorted := params.Sorted()
	values := sorted.InlineValues()
	return md5String(values + p.key)
}

func md5String(src string) string {
	hasher := md5.New()
	hasher.Write([]byte(src))
	return hex.EncodeToString(hasher.Sum(nil))
}
