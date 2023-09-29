package slice

import (
	"errors"
	"fmt"
)

func Delete[T any](slice []T, idx int) ([]T, error) {
	size := len(slice)
	if idx < 0 || idx >= size {
		return slice, errors.New(fmt.Sprintf("Invalid idx %d, slice valid size is %d", idx, size))
	}

	target := slice[:0]
	target = append(target, slice[:idx]...)
	target = append(target, slice[idx+1:]...)

	target = shrinkSlice(target)

	shrinkSlice(target)

	return target, nil
}

func shrinkSlice[T any](slice []T) []T {
	const threshold = 256
	oldCap := cap(slice)
	newCap := oldCap

	var triggerPoint int
	if oldCap <= threshold {
		// 小于256的时候 1/2 开始缩容，缩到 2/3
		newCap = oldCap * 2 / 3
		triggerPoint = oldCap / 2
	} else {
		// 大于256的时候 2/3 开始缩容，缩到 3/4
		newCap = oldCap * 3 / 4
		triggerPoint = oldCap * 2 / 3
	}

	// 8以下太小没有必要缩容
	if (len(slice)) == triggerPoint && triggerPoint >= 8 {
		newSlice := make([]T, triggerPoint, newCap)
		for index, value := range slice {
			newSlice[index] = value
		}

		return newSlice
	}

	return slice
}
