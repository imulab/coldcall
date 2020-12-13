package status

import (
	"github.com/imulab/coldcall"
	"net/http"
)

var (
	// Is returns a coldcall.Condition that evaluates to true when the status code matches.
	Is = func(status int) coldcall.Condition {
		return func(response *http.Response) bool {
			return response.StatusCode == status
		}
	}

	// InRange returns a coldcall.Condition that evaluates to true when the status is within
	// the supplied range.
	InRange = func(startInclusive int, endExclusive int) coldcall.Condition {
		return func(response *http.Response) bool {
			return response.StatusCode >= startInclusive && response.StatusCode < endExclusive
		}
	}

	// Is200 is a convenient wrapper for Is of 200 codes.
	Is200 = Is(http.StatusOK)

	// IsSuccess is a convenient wrapper for InRange for 2XX series codes.
	IsSuccess = InRange(200, 300)

	// IsFailure is a convenient wrapper for InRange for 4XX and 5XX series codes.
	IsFailure = InRange(400, 600)
)
