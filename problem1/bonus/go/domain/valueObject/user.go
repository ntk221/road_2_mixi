package valueObject

type UserID int

func NewUserID(id int) UserID {
	if id < 0 {
		panic("UserID must be positive")
	}

	return UserID(id)
}
