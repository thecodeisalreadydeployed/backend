package util

func MergeMaps(maps ...map[string]string) map[string]string {
	ret := make(map[string]string)
	for _, m := range maps {
		for k, v := range m {
			ret[k] = v
		}
	}
	return ret
}
