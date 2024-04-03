package util

func MapPtr[S any, T any](elements []*S, mpFunc func(*S) T) []T {

	var target = make([]T, 0, len(elements))
	for _, ele := range elements {
		target = append(target, mpFunc(ele))
	}
	return target
}

func Map[S any, T any](elements []S, mpFunc func(S) T) []T {

	var target = make([]T, 0, len(elements))
	for _, ele := range elements {
		target = append(target, mpFunc(ele))
	}
	return target
}

func CollectToMap[S any, K comparable, V any](elements []S, kvFunc func(S) (K, V)) map[K]V {

	var mp = make(map[K]V, len(elements))
	for _, ele := range elements {
		k, v := kvFunc(ele)
		mp[k] = v
	}
	return mp
}

// Difference compare src,dst justify which are new,and old
func Difference[T comparable](src, dst []T) ([]T, []T) {

	srcMap := CollectToMap[T, T, struct{}](src, func(t T) (T, struct{}) { return t, struct{}{} })
	dstMap := CollectToMap[T, T, struct{}](dst, func(t T) (T, struct{}) { return t, struct{}{} })

	var newItemList, oldItemList = []T{}, []T{}

	for _, ele := range src {
		if _, ok := dstMap[ele]; !ok { //new
			newItemList = append(newItemList, ele)
		}
	}

	for _, ele := range dst {
		if _, ok := srcMap[ele]; !ok { //old
			oldItemList = append(oldItemList, ele)
		}

	}
	return newItemList, oldItemList
}
