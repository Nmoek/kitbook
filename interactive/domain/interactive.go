package domain

type Interactive struct {
	BizId      int64
	Biz        string
	ReadCnt    int64
	LikeCnt    int64
	CollectCnt int64
	Liked      bool
	Collected  bool
}
