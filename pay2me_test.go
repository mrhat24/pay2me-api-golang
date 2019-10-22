package pay2me_api

import (
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
