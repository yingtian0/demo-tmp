package domain

import (
	"strings"
	"time"
)

type User struct {
	id          UserID
	displayName string
	profileText string
	createdAt   time.Time
}

func NewUser(
	id UserID,
	displayName string,
	profileText string,
	now time.Time,
) (*User, error) {
	if id.IsZero() || isBlank(displayName) || now.IsZero() {
		return nil, ErrInvalidUser
	}

	return &User{
		id:          id,
		displayName: strings.TrimSpace(displayName),
		profileText: strings.TrimSpace(profileText),
		createdAt:   now,
	}, nil
}

func (u *User) ID() UserID {
	return u.id
}

func (u *User) DisplayName() string {
	return u.displayName
}

func (u *User) ProfileText() string {
	return u.profileText
}

func (u *User) CreatedAt() time.Time {
	return u.createdAt
}
