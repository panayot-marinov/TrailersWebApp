package src

import "time"

// each session contains the username of the user and the time at which it expires
type Session struct {
	username string
	expiry   time.Time
}

// we'll use this method later to determine if the session has expired
func (s Session) isExpired() bool {
	return s.expiry.Before(time.Now())
}
