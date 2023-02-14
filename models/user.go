package models

type User struct {
	UserID   int64  `db:"user_id"`
	Username string `db:"username"`
	Password string `db:"password"`
	FollowCount   int64  `db:"follow_count"`
	FollowerCount int64  `db:"follower_count"`
	IsFollow      bool   `db:"is_follow"`
}

