package enum

type DailyLessonStatusType int32

const (
	UnknownDailyLessonStatusType DailyLessonStatusType = 20040 + iota
	DailyLessonStatusType_PENDING
	DailyLessonStatusType_START
	DailyLessonStatusType_COMPLETED
	DailyLessonStatusType_CANCELED
)

var (
	// 未开始--（开始上课）-->上课中--（点名）(结束上课) ->已结束--(上传上课记录)
	DailyLessonStatusType_name = map[int32]string{
		20040: "未知来源",
		20041: "未开始",
		20042: "上课中",
		20043: "已结束",
		20044: "已取消",
	}
	DailyLessonStatusType_value = map[string]int32{}
)

func (n DailyLessonStatusType) String() string {
	if name, exist := DailyLessonStatusType_name[int32(n)-20040]; exist {
		return name
	}
	return ""
}

type DailyLessonStudentStatusType int32

const (
	UnknownDailyLessonStudentStatusTypeType DailyLessonStudentStatusType = 20050 + iota
	DailyLessonStudentStatusType_UNSIGNED
	DailyLessonStudentStatusType_SIGNED
	DailyLessonStudentStatusType_ABSENT
	DailyLessonStudentStatusType_CANCELED
)

var (
	DailyLessonStudentStatusType_name = map[int32]string{
		20050: "未知来源",
		20051: "未签到",
		20052: "已签到",
		20053: "请假",
		20054: "已取消",
	}
	DailyLessonStudentStatusType_value = map[string]int32{}
)

func (n DailyLessonStudentStatusType) String() string {
	if name, exist := DailyLessonStudentStatusType_name[int32(n)-20050]; exist {
		return name
	}
	return ""
}
