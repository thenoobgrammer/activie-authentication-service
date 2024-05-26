package redis

import "fmt"

var (
	ErrKeyNotFound              = fmt.Errorf("key not found")
	ErrKeyFormatInvalid         = fmt.Errorf("key format invalid")
	ErrScanningKeys             = fmt.Errorf("error scanning keys")
	ErrEmailNotFound            = fmt.Errorf("no email found for token")
	ErrGettingToken             = fmt.Errorf("error getting token")
	ErrSettingToken             = fmt.Errorf("error setting token")
	ErrSettingExpirationOnToken = fmt.Errorf("error setting expiration time on token")
	ErrPersitingToken           = fmt.Errorf("error invalidating expiration on current token")
)
