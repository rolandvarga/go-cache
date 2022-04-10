package main

import (
	"reflect"
	"testing"
)

func testServiceWith(entries []Entry) *Service {
	s := newService()
	s.cache.Entries = entries
	return s
}

func TestCacheAddEntry(t *testing.T) {
	t.Parallel()
	cases := map[string]struct {
		entries []Entry
		this    Entry
		want    []Entry
	}{
		"Add to empty": {
			entries: []Entry{},
			this:    Entry{UUID: "1234", TTL: 10, Body: "asdf"},
			want: []Entry{
				Entry{UUID: "1234", TTL: 10, Body: "asdf"},
			},
		},
		"Add to multiple": {
			entries: []Entry{
				Entry{UUID: "1234", TTL: 10, Body: "asdf"},
				Entry{UUID: "5678", TTL: 20, Body: "qwer"},
			},
			this: Entry{UUID: "9012", TTL: 30, Body: "zxcv"},
			want: []Entry{
				Entry{UUID: "1234", TTL: 10, Body: "asdf"},
				Entry{UUID: "5678", TTL: 20, Body: "qwer"},
				Entry{UUID: "9012", TTL: 30, Body: "zxcv"},
			},
		},
	}

	for name, c := range cases {
		c := c
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			s := testServiceWith(c.entries)

			s.cache.add(c.this)
			if !reflect.DeepEqual(s.cache.Entries, c.want) {
				t.Errorf("have '%v' want '%v'", s.cache.Entries, c.want)
			}
		})
	}
}

func TestCacheRemoveEntry(t *testing.T) {
	t.Parallel()
	cases := map[string]struct {
		entries []Entry
		this    Entry
		want    []Entry
	}{
		"Element in middle": {
			entries: []Entry{
				Entry{UUID: "1234", TTL: 10, Body: "asdf"},
				Entry{UUID: "5678", TTL: 20, Body: "qwer"},
				Entry{UUID: "9012", TTL: 30, Body: "zxcv"},
			},
			this: Entry{UUID: "5678", TTL: 20, Body: "qwer"},
			want: []Entry{
				Entry{UUID: "1234", TTL: 10, Body: "asdf"},
				Entry{UUID: "9012", TTL: 30, Body: "zxcv"},
			},
		},
		"First element": {
			entries: []Entry{
				Entry{UUID: "1234", TTL: 10, Body: "asdf"},
				Entry{UUID: "5678", TTL: 20, Body: "qwer"},
				Entry{UUID: "9012", TTL: 30, Body: "zxcv"},
			},
			this: Entry{UUID: "1234", TTL: 10, Body: "asdf"},
			want: []Entry{
				Entry{UUID: "5678", TTL: 20, Body: "qwer"},
				Entry{UUID: "9012", TTL: 30, Body: "zxcv"},
			},
		},
		"Last element odd": {
			entries: []Entry{
				Entry{UUID: "1234", TTL: 10, Body: "asdf"},
				Entry{UUID: "5678", TTL: 20, Body: "qwer"},
				Entry{UUID: "9012", TTL: 30, Body: "zxcv"},
			},
			this: Entry{UUID: "9012", TTL: 30, Body: "zxcv"},
			want: []Entry{
				Entry{UUID: "1234", TTL: 10, Body: "asdf"},
				Entry{UUID: "5678", TTL: 20, Body: "qwer"},
			},
		},
		"Last element even": {
			entries: []Entry{
				Entry{UUID: "1234", TTL: 10, Body: "asdf"},
				Entry{UUID: "5678", TTL: 20, Body: "qwer"},
				Entry{UUID: "9012", TTL: 30, Body: "zxcv"},
				Entry{UUID: "3456", TTL: 40, Body: "sdfg"},
			},
			this: Entry{UUID: "3456", TTL: 40, Body: "sdfg"},
			want: []Entry{
				Entry{UUID: "1234", TTL: 10, Body: "asdf"},
				Entry{UUID: "5678", TTL: 20, Body: "qwer"},
				Entry{UUID: "9012", TTL: 30, Body: "zxcv"},
			},
		},
	}

	for name, c := range cases {
		c := c
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			s := testServiceWith(c.entries)

			s.cache.remove(c.this)
			if !reflect.DeepEqual(s.cache.Entries, c.want) {
				t.Errorf("have '%v' want '%v'", s.cache.Entries, c.want)
			}
		})
	}

}
