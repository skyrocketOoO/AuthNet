package domain

type ErrGraphCycle struct{}

func (e ErrGraphCycle) Error() string {
	return "cycle detected"
}

type ErrBodyAttribute struct{}

func (e ErrBodyAttribute) Error() string {
	return "body attribute error"
}

type ErrNotImplemented struct{}

func (e ErrNotImplemented) Error() string {
	return "not implemented"
}

type ErrRecordNotFound struct{}

func (e ErrRecordNotFound) Error() string {
	return "record not found"
}

type ErrDuplicateRecord struct{}

func (e ErrDuplicateRecord) Error() string {
	return "duplicate record"
}
