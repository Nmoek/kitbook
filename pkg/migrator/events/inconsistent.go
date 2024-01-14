package events

type InconsistentEvent struct {
	Id        int64
	Type      string
	Direction string
}

const (
	//校验的目标数据，缺了这一条
	InconsistentEventTypeTargetMissing = "target_missing"
	//数据不相等
	InconsistentEventTypeNEQ = "neq"
	//校验的源数据，缺了这一条
	InconsistentEventTypeBaseMissing = "base_missing"
)
