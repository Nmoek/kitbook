package domain

type Payment struct {
	Amt Amount

	BizTradeNO  string
	Description string
	Status      PaymentStatus
	TxnID       string //第三方返回的ID
}

type Amount struct {
	Total    int64
	Currency string
}

type PaymentStatus uint8

func (p PaymentStatus) AstoUint8() uint8 { return uint8(p) }

const (
	PaymentStatusUnknown = iota
	// 扫码但还未支付
	PaymentStatusInit
	// 支付成功
	PaymentStatusSuccess
	// 支付失败
	PaymentStatusFail
	// 退款中
	PaymentStatusRefund
)
