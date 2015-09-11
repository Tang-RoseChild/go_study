package shorturl

import (
	"testing"
)

func TestGet(t *testing.T) {
	urlstore := &URLStore{
		urls: map[string]string{"a": "baidu.com", "b": "hao123.com", "c": "fenghuang.com"},
	}

	testData := []struct{ in, out string }{{"a", "baidu.com"}, {"b", "hao123.com"}, {"d", ""}, {"e", "google.com"}}
	for i, tt := range testData {
		s := urlstore.Get(tt.in)
		if s != tt.out {
			t.Errorf("%d. %q => %q, wanted : %q \n ", i, tt.in, s, tt.out)
		}
	}
}

func TestSet(t *testing.T) {
	urlstore := &URLStore{
		urls: make(map[string]string),
	}

	testData := []struct{ key, value string }{{"a", "baidu.com"}, {"b", "hao123.com"}, {"d", ""}, {"a", "google.com"}}
	for i, tt := range testData {
		r := urlstore.Set(tt.key, tt.value)

		if i == 3 && r {
			t.Fatal("can not save the same key more than once")
		}

		if i != 3 && !r {
			t.Fatal("set error  ", i, " : ", tt.key, " : ", tt.value)
		}

	}

}

func TestCount(t *testing.T) {
	urlstore := &URLStore{
		urls: map[string]string{"a": "baidu.com", "b": "hao123.com", "c": "fenghuang.com"},
	}

	stores := []*URLStore{urlstore, &URLStore{
		urls: map[string]string{},
	}}

	for i, store := range stores {
		c := store.Count()
		// c = 7 // for testing fail
		if c != len(store.urls) {
			t.Fatalf("%d Count => %d, wanted %d \n", i, c, len(urlstore.urls))
		}
	}
}
