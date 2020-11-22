package utils

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func Partition(sliceSize int, partitionSize int, do func(int, int) error) error {
	for i := 0; i < sliceSize; i += partitionSize {
		err := do(i, Min(i+partitionSize, sliceSize))
		if err != nil {
			return err
		}
	}
	return nil
}
