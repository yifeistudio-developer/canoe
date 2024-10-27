package ports

import "github.com/yifeistudio-developer/canoe/internal/application/core/domain"

type DbPort interface {
	UserDbPort() UserDbPort
}

type UserDbPort interface {
	GetById(id int64) (domain.User, error)
	Save(user domain.User) error
}

type SessionDbPort interface {
	GetById(id int64) (domain.Session, error)
	Save(session domain.Session) error
}

type UserSessionDbPort interface {
	GetById(id int64) (domain.UserSession, error)
	Save(session domain.UserSession) error
}

type GroupDbPort interface {
	GetById(id int64) (domain.Group, error)
	Save(group domain.Group) error
}

type GroupMemberDbPort interface {
	GetById(id int64) (domain.GroupMember, error)
	Save(group domain.GroupMember) error
}

type MessageDbPort interface {
}

type EventDbPort interface {
}
