package enum

type CategoryType int32

const (
	CategoryType_Unknown EnableStatus = 10020 + iota
	CategoryType_PAINTING
	CategoryType_CALLIGRAPHY
)

var (
	CategoryType_name = map[int32]string{
		10020: "未知",
		10021: "绘画",
		10022: "书法",
	}
	CategoryType_value = map[string]int32{}
)

func (n CategoryType) String() string {
	if name, exist := CategoryType_name[int32(n)-10020]; exist {
		return name
	}
	return ""
}
