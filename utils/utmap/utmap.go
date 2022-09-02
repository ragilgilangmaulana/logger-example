package utmap

// Merge function
//
// **@Params:** [ `a`: map 1; `b`: map 2 ]
func Merge(a *map[string]interface{}, b map[string]interface{}) {
	for k, v := range b {
		if _, ok := b[k]; ok {
			(*a)[k] = v
		}
	}
}
