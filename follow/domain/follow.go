package domain

type FollowRelation struct {
	Followee int64
	Follower int64
}

type FollowStatics struct {
	Followees int64 // 关注数
	Followers int64 // 粉丝数
}
