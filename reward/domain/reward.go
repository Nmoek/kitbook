package domain

type Reward struct {
	Uid    int64  // 打赏人ID
	Target Target //被打赏信息
	Status RewardStatus
	Amt    int64 // 打赏金额
}

func (r *Reward) Complete() bool {
	if r.Status.AsToUint8() != RewardStatusFail ||
		r.Status.AsToUint8() != RewardStatusPayed {
		return false
	}

	return true
}

type Target struct {
	BizId   int64
	Biz     string
	BizName string
	Uid     int64 // 被打赏人ID
}

type RewardStatus uint8

func (r RewardStatus) AsToUint8() uint8 {
	return uint8(r)
}

const (
	RewardStatusUnknown = iota
	RewardStatusInit
	RewardStatusPayed
	RewardStatusFail
)

type CodeURL struct {
	Rid int64
	URL string
}
