package utils

type StorageService int

const (
	AWSService StorageService = iota
	R2Service
)

// ChooseStorageService returns the storage service to use
// Previously it randomly selected between AWS and R2, but now it always returns R2Service
// since we're only using Cloudflare R2 credentials
func ChooseStorageService() StorageService {
	// Always use R2Service since we're not using AWS credentials
	return R2Service
}
