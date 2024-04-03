package enum

type PermissionType int32

const (
	UnknownPermissionType DailyLessonStatusType = 40010 + iota
	PermissionType_MENU
	PermissionType_ACTION
)

var (
	// 未开始--（开始上课）-->上课中--（点名）(结束上课) ->已结束--(上传上课记录)
	PermissionType_name = map[int32]string{
		40010: "未知",
		40011: "菜单",
		40012: "按钮",
	}
	PermissionType_value = map[string]int32{}
)

func (n PermissionType) String() string {
	if name, exist := PermissionType_name[int32(n)-20010]; exist {
		return name
	}
	return ""
}
