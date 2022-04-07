package main

import (
	"reflect"
	"testing"
)

// TODO []Entry
func removeFromSlice(slice []int, idx int) []int {
	return append(slice[:idx], slice[idx+1:]...)
}

func TestRemoveFromSlice(t *testing.T) {
	t.Parallel()
	cases := map[string]struct {
		slice []int
		idx   int
		want  []int
	}{
		"Element in middle": {
			slice: []int{1, 2, 3},
			idx:   1,
			want:  []int{1, 3},
		},
		"First element": {
			slice: []int{1, 2, 3},
			idx:   0,
			want:  []int{2, 3},
		},
		"Last element odd": {
			slice: []int{1, 2, 3},
			idx:   2,
			want:  []int{1, 2},
		},
		"Last element even": {
			slice: []int{1, 2, 3, 4},
			idx:   3,
			want:  []int{1, 2, 3},
		},
	}

	for name, c := range cases {
		c := c
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := removeFromSlice(c.slice, c.idx)
			if !reflect.DeepEqual(got, c.want) {
				t.Errorf("got '%v' want '%v'", got, c.want)
			}
		})
	}

}
