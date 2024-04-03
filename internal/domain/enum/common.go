package enum

type EnableStatus int32

const (
	EnableStatusUnknown EnableStatus = 10010 + iota
	EnableStatusEnabled
	EnableStatusDiaabled
)

var (
	EnableStatus_name = map[int32]string{
		10010: "未知",
		10011: "启用",
		10012: "禁用",
	}
	EnableStatus_value = map[string]int32{}
)

func (n EnableStatus) String() string {
	if name, exist := EnableStatus_name[int32(n)-10010]; exist {
		return name
	}
	return ""
}

type MetaDataType int32

const (
	CourseCodeSequence MetaDataType = 90010 + iota
)

var (
	MetaData_name = map[int32]string{
		90010: "课程编码序列号",
	}
	MetaData_value = map[string]int32{}
)

func (n MetaDataType) String() string {
	if name, exist := MetaData_name[int32(n)-90010]; exist {
		return name
	}
	return ""
}
