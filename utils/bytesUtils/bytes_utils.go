package byteUtils

func TotalInArray(bytesArray [][]byte) int {
	total := 0

	for _, bytes := range bytesArray {
		total += len(bytes)
	}

	return total
}
