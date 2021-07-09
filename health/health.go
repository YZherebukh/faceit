// ------------go:generate mockgen -source ../health/health.go -destination ../health/mock/mock_health.go

package health

import "time"

// Response is a struct with health response definition
type Response struct {
	Name    string    `json:"name"`
	Healthy bool      `json:"healthy"`
	Time    time.Time `json:"time"`
	Message string    `json:"message"`
}
