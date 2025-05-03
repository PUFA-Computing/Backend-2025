package utils

// No imports needed

type StorageService int

const (
	AWSService StorageService = iota
	R2Service
)

// ChooseStorageService always returns R2Service as we no longer use AWS S3
func ChooseStorageService() StorageService {
	// Always use R2 service
	return R2Service
}
