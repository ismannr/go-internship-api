package models

type Role string

const (
	RoleMentor      Role = "mentor"
	RoleParticipant Role = "participant"
)

type Level string

const (
	LevelAdmin Level = "admin"
	LevelUser  Level = "user"
)
