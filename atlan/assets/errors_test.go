package assets

import (
	"errors"
	"strings"
	"testing"
)

// TestAtlanError_PassthroughNoFormatNoise reproduces the kernel-ext-atlan "530 /
// dead-DNS tenant" failure: an HTTP passthrough error whose real detail lives in
// OriginalError. Error() must not emit %!(MISSING)/%!(EXTRA) format noise, and
// must still surface the original server response.
func TestAtlanError_PassthroughNoFormatNoise(t *testing.T) {
	orig := errors.New(`API returned status code 530: {"error_code":1016,"detail":"origin DNS error"}`)
	// handleApiError routes a 5xx to ERROR_PASSTHROUGH with a single (often empty)
	// causes string — fewer args than the old template's verbs.
	err := ThrowAtlanError(orig, ERROR_PASSTHROUGH, nil, "")
	msg := err.Error()

	if strings.Contains(msg, "%!") {
		t.Fatalf("Error() contains format noise: %q", msg)
	}
	if !strings.Contains(msg, "status code 530") || !strings.Contains(msg, "1016") {
		t.Fatalf("Error() dropped the original server response: %q", msg)
	}
	if !strings.Contains(msg, "ATLAN-GO-500-000") {
		t.Fatalf("Error() missing error id: %q", msg)
	}
}

// TestAtlanError_Format covers the Error() substitution rules directly.
func TestAtlanError_Format(t *testing.T) {
	cases := []struct {
		name    string
		err     AtlanError
		want    string // substring that must be present
		notWant string // substring that must be absent ("" = skip)
	}{
		{
			name: "verb template with matching arg formats once",
			err:  AtlanError{ErrorCode: ErrorInfo{ErrorID: "X", ErrorMessage: "category: %s."}, Args: []interface{}{"myCategory"}},
			want: "category: myCategory.",
		},
		{
			name:    "verb template with no args does not produce MISSING noise",
			err:     AtlanError{ErrorCode: ErrorInfo{ErrorID: "X", ErrorMessage: "code: %s, %s"}},
			notWant: "%!",
		},
		{
			name:    "static template with stray args does not produce EXTRA noise",
			err:     AtlanError{ErrorCode: ErrorInfo{ErrorID: "X", ErrorMessage: "static message"}, Args: []interface{}{"ignored"}},
			notWant: "%!",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			msg := tc.err.Error()
			if tc.want != "" && !strings.Contains(msg, tc.want) {
				t.Fatalf("want %q in %q", tc.want, msg)
			}
			if tc.notWant != "" && strings.Contains(msg, tc.notWant) {
				t.Fatalf("did not want %q in %q", tc.notWant, msg)
			}
		})
	}
}
