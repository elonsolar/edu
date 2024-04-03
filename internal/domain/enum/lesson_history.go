package enum

type NumChangeSourceType int32

const (
	UnknownSourceType NumChangeSourceType = 20010 + iota
	ManualAdjust
	Order
	RollCall
)

var (
	NumChangeSourceType_name = map[int32]string{
		20010: "未知来源",
		20011: "手工调整",
		20012: "订单购买",
		20013: "点名签到",
	}
	ErrorReason_value = map[string]int32{}
)

func (n NumChangeSourceType) String() string {
	if name, exist := NumChangeSourceType_name[int32(n)-20010]; exist {
		return name
	}
	return ""
}
