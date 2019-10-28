package pay2me_api

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"io"
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

type Deals []Deal

func (d *Deals) CompleteJson() io.Reader {
	var dealsIds []string
	for _, d := range *d {
		dealsIds = append(dealsIds, d.ObjectID)
	}
	dealsMap := map[string][]string{
		"deals": dealsIds,
	}
	j, _ := json.Marshal(dealsMap)
	return bytes.NewReader(j)
}

// required fields: OrderAmount, OrderDesc, ObjectID
func (d *Deal) CreationJSON() Pay2MeParams {
	return Pay2MeParams{
		"order_amount": d.OrderAmount,
		"order_desc":   d.OrderDesc,
		"order_id":     d.ObjectID,
	}
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

func (params Pay2MeParams) Json(key string) io.Reader {
	p := params.Sorted()
	p["signature"] = GetSignature(p, key)
	j, _ := json.Marshal(p)
	return bytes.NewReader(j)
}

func (params Pay2MeParams) InlineValues() string {
	result := ""
	for k := range params {
		result += params[k]
	}
	return result
}

type Pay2MeApi struct {
	key    string
	ApiUrl string
}

func (p *Pay2MeApi) doRequest(r *http.Request) (*http.Response, error) {
	c := http.DefaultClient
	r.Header.Add("X-API-KEY", p.key)
	r.Header.Add("Accept", "application/json")
	return c.Do(r)
}

func GetSignature(params Pay2MeParams, key string) string {
	sorted := params.Sorted()
	values := sorted.InlineValues()
	return md5String(values + key)
}

func (p *Pay2MeApi) DealCreate(deal *Deal) (*http.Response, error) {
	dealParams := deal.CreationJSON()
	r, err := http.NewRequest("POST", p.ApiUrl+"/deals", dealParams.Json(p.key))
	if err != nil {
		return nil, err
	}
	return p.doRequest(r)
}

func (p *Pay2MeApi) DealStatus(deal *Deal) (*http.Response, error) {
	r, err := http.NewRequest("GET", p.ApiUrl+"/deals/status/"+deal.ObjectID, nil)
	if err != nil {
		return nil, err
	}
	return p.doRequest(r)
}

func (p *Pay2MeApi) DealComplete(deal *Deal) (*http.Response, error) {
	r, err := http.NewRequest("PUT", p.ApiUrl+"/deals/complete/"+deal.ObjectID, nil)
	if err != nil {
		return nil, err
	}
	return p.doRequest(r)
}

func (p *Pay2MeApi) DealsComplete(deals *Deals) (*http.Response, error) {
	r, err := http.NewRequest("PUT", p.ApiUrl+"/deals/complete", deals.CompleteJson())
	if err != nil {
		return nil, err
	}
	return p.doRequest(r)
}

func (p *Pay2MeApi) DealCancel(deal *Deal) (*http.Response, error) {
	r, err := http.NewRequest("PUT", p.ApiUrl+"/deals/cancel/"+deal.ObjectID, nil)
	if err != nil {
		return nil, err
	}
	return p.doRequest(r)
}

func CreatePay2MeApi(key string) *Pay2MeApi {
	api := &Pay2MeApi{
		key:    key,
		ApiUrl: "https://api.pay2me.world/api/v3",
	}
	return api
}

func md5String(src string) string {
	hasher := md5.New()
	hasher.Write([]byte(src))
	return hex.EncodeToString(hasher.Sum(nil))
}
