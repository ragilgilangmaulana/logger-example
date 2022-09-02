package uttype

import "fmt"

// TypeToString function
//
// **@Params:** [ `v`: any ]
//
// **@Returns:** [ `$1`: string type ]
func TypeToString(v interface{}) string {
	return fmt.Sprintf("%T", v)
}
