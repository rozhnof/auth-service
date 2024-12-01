package repo

import "github.com/google/uuid"

type UserFilters struct {
	UserIDs []uuid.UUID
}

type Pagination struct {
	Limit  int32
	Offset int32
}
