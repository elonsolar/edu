package enum

type CoursePlanExcludeDateType int32

const (
	ExcludeDateTypeUnknown NumChangeSourceType = 30010 + iota
	CoursePlanExcludeDateType_WeekDay
	CoursePlanExcludeDateType_TimeInterval
)

var (
	CoursePlanExcludeDateType_name = map[int32]string{
		30010: "未知",
		30011: "一周中的一天",
		30012: "时间段",
	}
	CoursePlanExcludeDateType_value = map[string]int32{}
)

func (n CoursePlanExcludeDateType) String() string {
	if name, exist := CoursePlanExcludeDateType_name[int32(n)-30010]; exist {
		return name
	}
	return ""
}

type CoursePlanStatusType int32

const (
	CoursePlanStatusTypeUnknown NumChangeSourceType = 30020 + iota
	CoursePlanStatusType_NEW
	CoursePlanStatusType_RELEASED
	CoursePlanStatusType_SCHEDULED
	CoursePlanStatusType_CANCELED
)

var (
	CoursePlanStatusType_name = map[int32]string{
		30020: "未知",
		30021: "新纪录",
		30022: "已发布",
		30023: "已排期",
		30024: "已取消",
	}
	CoursePlanStatusType_value = map[string]int32{}
)

func (n CoursePlanStatusType) String() string {
	if name, exist := CoursePlanExcludeDateType_name[int32(n)-30020]; exist {
		return name
	}
	return ""
}

type CoursePlanDetailStatusType int32

const (
	CoursePlanDetailStatusTypeUnknown NumChangeSourceType = 30030 + iota
	CoursePlanDetailStatusType_NEW
	CoursePlanDetailStatusType_SCHEDULED
	CoursePlanDetailStatusType_STOPPED
)

var (
	CoursePlanDetailStatusType_name = map[int32]string{
		30030: "未知",
		30031: "新纪录",
		30032: "已排期",
		30033: "已停课",
	}
	CoursePlanDetailStatusType_value = map[string]int32{}
)

func (n CoursePlanDetailStatusType) String() string {
	if name, exist := CoursePlanExcludeDateType_name[int32(n)-30030]; exist {
		return name
	}
	return ""
}

type CoursePlanStudentStatusType int32

const (
	CoursePlanStudentStatusType_UNKNOWN NumChangeSourceType = 30040 + iota
	CoursePlanStudentStatusType_NEW
	CoursePlanStudentStatusType_SCHEDULED
	CoursePlanStudentStatusType_STOPPED
)

var (
	CoursePlanStudentStatusType_name = map[int32]string{
		30040: "未知",
		30041: "新纪录",
		30042: "已排期",
		30043: "已停课",
	}
	CoursePlanStudentStatusType_value = map[string]int32{}
)

func (n CoursePlanStudentStatusType) String() string {
	if name, exist := CoursePlanStudentStatusType_name[int32(n)-30040]; exist {
		return name
	}
	return ""
}
