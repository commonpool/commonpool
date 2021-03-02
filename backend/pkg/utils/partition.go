package utils

func Partition(sliceSize int, partitionSize int, do func(int, int) error) error {
	for i := 0; i < sliceSize; i += partitionSize {
		err := do(i, Min(i+partitionSize-1, sliceSize-1))
		if err != nil {
			return err
		}
	}
	return nil
}
