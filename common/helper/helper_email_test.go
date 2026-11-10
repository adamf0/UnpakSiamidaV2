package helper

import "testing"

func TestValidateOptionalEmail(t *testing.T) {
	tests := []struct {
		name    string
		email   interface{}
		wantErr bool
	}{
		{name: "empty string allowed", email: "", wantErr: false},
		{name: "whitespace allowed", email: "   ", wantErr: false},
		{name: "nil allowed", email: nil, wantErr: false},
		{name: "valid gmail allowed", email: "user@gmail.com", wantErr: false},
		{name: "valid unpak allowed", email: "admin@unpak.ac.id", wantErr: false},
		{name: "invalid format rejected", email: "not-an-email", wantErr: true},
		{name: "double dot rejected", email: "user..name@unpak.ac.id", wantErr: true},
		{name: "space inside email rejected", email: "user name@unpak.ac.id", wantErr: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateOptionalEmail(tc.email)
			if (err != nil) != tc.wantErr {
				t.Fatalf("ValidateOptionalEmail(%v) err = %v, wantErr = %v", tc.email, err, tc.wantErr)
			}
		})
	}
}
