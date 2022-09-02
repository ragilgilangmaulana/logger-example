package utmap

import (
	"github.com/koinworks/asgard-heimdal/utils/utstring"
)

func MergeString(a *map[string]string, b map[string]string) {
	for k, v := range b {
		if _, ok := b[k]; ok {
			(*a)[k] = v
		}
	}
}

// MatrixMapString function
//
// **@Params:** [ `a`: map 1; `b`: map 2 ]
//
// **@Returns:** [ `$1`: matrix map ]
func MatrixMapString(s []map[string]string, d []map[string]string) []map[string]string {
	newMaps := []map[string]string{}
	for _, v := range s {
		for _, v2 := range d {
			cur := map[string]string{}
			utstring.MergeString(&cur, v)
			utstring.MergeString(&cur, v2)
			newMaps = append(newMaps, cur)
		}
	}
	return newMaps
}
