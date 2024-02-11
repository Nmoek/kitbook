package events

// PaymentEvent
// @Description: 后续需要接入大数据业务时再进行加字段
type PaymentEvent struct {
	BizTradeNO string
	Status     uint8
}

func (p PaymentEvent) Topic() string {
	return "payment_events"
}
