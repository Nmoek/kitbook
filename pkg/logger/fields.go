package logger

type Field struct {
	Key string
	Val any
}

func Error(err error) Field {
	return Field{
		Key: "error",
		Val: err,
	}
}

func Int[T intVal](key string, val T) Field {
	return Field{
		Key: key,
		Val: val,
	}
}

func Float[T floatVal](key string, val T) Field {
	return Field{
		Key: key,
		Val: val,
	}
}

type intVal interface {
	int | int8 | int16 | int32 | int64 |
		uint | uint8 | uint16 | uint32 | uint64
}

type floatVal interface {
	float32 | float64
}

type Fields []Field

// @func: Add
// @date: 2023-11-23 00:20:05
// @brief: 添加Field
// @author: Kewin Li
// @receiver fs
// @param feild
// @return *Fields
func (fs Fields) Add(feild Field) Fields {
	return append(fs, feild)
}
