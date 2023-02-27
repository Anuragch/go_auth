package errors

type DBError struct {
	S string
}

func (err *DBError) Error() string {
	return err.S
}
