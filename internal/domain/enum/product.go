package enum

type SkuCategoryType int32

const (
	SkuCategoryType_UNKNOWN SkuCategoryType = 20070 + iota
	SkuCategoryType_LESSON
	SkuCategoryType_BOOK
	SkuCategoryType_OTHER
)

var (
	SkuCategoryType_name = map[int32]string{
		20070: "未知",
		20071: "课时",
		20072: "教材",
		20073: "其他",
	}
	SkuCategoryType_value = map[string]int32{}
)

func (n SkuCategoryType) String() string {
	if name, exist := SkuCategoryType_name[int32(n)-20070]; exist {
		return name
	}
	return ""
}

type SkuStatusType int32

const (
	SkuStatusType_UNKNOWN SkuStatusType = 20080 + iota
	SkuStatusType_AVAILABLE
	SkuStatusType_DISCONTINUED
	SkuStatusType_OTHER
)

var (
	SkuStatusType_name = map[int32]string{
		20080: "未知",
		20081: "在售",
		20082: "下架",
	}

	SkuStatusType_value = map[string]int32{}
)

func (n SkuStatusType) String() string {
	if name, exist := SkuStatusType_name[int32(n)-20080]; exist {
		return name
	}
	return ""
}

type SkuType int32

const (
	SkuType_UNKNOWN SkuStatusType = 20090 + iota
	SkuType_SINGLE
	SkuType_COMBINE
)

var (
	SkuType_name = map[int32]string{
		20090: "未知",
		20091: "单品",
		20092: "套餐",
	}

	SkuType_value = map[string]int32{}
)

func (n SkuType) String() string {
	if name, exist := SkuType_name[int32(n)-20090]; exist {
		return name
	}
	return ""
}
