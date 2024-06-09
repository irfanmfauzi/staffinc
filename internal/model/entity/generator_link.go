package entity

import "time"

type GeneratorLink struct {
	Id          int64     `db:"id"`
	UserId      int64     `db:"user_id"`
	Code        string    `db:"code"`
	ExpiredAt   time.Time `db:"expired_at"`
	CountAccess int       `db:"count_access"`
}
