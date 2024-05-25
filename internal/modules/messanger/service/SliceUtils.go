package service

func ReverseStringSlice(slice []string) []string {
	for i := 0; i < len(slice)/2; i++ {
		tmp := slice[i]
		slice[i] = slice[len(slice)-i-1]
		slice[len(slice)-i-1] = tmp
	}
	return slice
}
