package value

type ErrCode uint8

const (
	ErrCodeAuthorization         = 10
	ErrCodeRecordNotExist        = 11
	ErrCodeOther                 = 12
	ErrCodeDatabaseTimeout       = 13
	ErrCodeInvalidParameter      = 14
	ErrCodeNoResultFound         = 15
	ErrCodeDuplicate             = 16
	ErrCodePostalCodeNotFound    = 17
	ErrCodeDailyLimitExceeded    = 18
	ErrCodeHourlyLimitExceeded   = 19
	ErrCodeWeeklyLimitExceeded   = 20
	ErrCodeInvalidInput          = 21
	ErrCodeServerOverloaded      = 22
	ErrCodeServiceNotImplemented = 23
	ErrCodeRadiusTooLarge        = 24
	ErrCodeMaxRowsTooLarge       = 27
)
