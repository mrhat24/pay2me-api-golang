package pay2me_api

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"testing"
)

func TestPay2MeParams_Sort(t *testing.T) {
	m := Pay2MeParams{
		"foo": "foo",
		"bar": "bar",
		"biz": "biz",
	}
	newMap := m.Sorted()
	result := ""
	for k := range newMap {
		result += k
	}
	if result != "barbizfoo" {
		t.Errorf("sorting are wrong")
	}
}

func TestPay2MeParams_InlineValues(t *testing.T) {
	m := Pay2MeParams{
		"foo": "1",
		"bar": "2",
		"biz": "3",
	}
	m = m.Sorted()
	result := m.InlineValues()
	if result != "231" {
		t.Errorf("values wrong")
	}
}

func TestPay2MeParams_Json(t *testing.T) {
	m := Pay2MeParams{
		"goo": "4",
		"foo": "1",
		"bar": "2",
		"biz": "3",
	}
	j := m.Json("key")
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(j)
	if err != nil {
		t.Error(err)
	}
	s := buf.String()
	if len(s) == 0 {
		t.Errorf("len is 0")
	}
}

func TestPay2MeApi_DealCreate(t *testing.T) {
	p2m := CreatePay2MeApi("")
	r, _ := p2m.DealCreate(&Deal{
		CreateDate:  "",
		OrderID:     "",
		UpdateDate:  "",
		ObjectID:    "",
		Redirect:    "",
		OrderAmount: "",
		Signature:   "",
		ExpireDate:  "",
		OrderDesc:   "",
		Status:      "",
	})
	body, _ := ioutil.ReadAll(r.Body)
	fmt.Println(string(body))
}
