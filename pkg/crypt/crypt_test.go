//go:build unit

package crypt

import "testing"

func Test_encrypt(t *testing.T) {
	type args struct {
		plaintext string
		key       string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "crypt/encrypt test",
			args: args{
				plaintext: "interview_test_pass",
				key:       "04076314bdb5ecf31706eea86fc984d6"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ciphertext, err := Encrypt(tt.args.plaintext, tt.args.key)
			if err != nil {
				t.Errorf("encrypt() error = %v", err)
				return
			}
			plaintext, err := Decrypt(ciphertext, tt.args.key)
			if err != nil {
				t.Errorf("encrypt() error = %v", err)
				return
			}
			if plaintext != tt.args.plaintext {
				t.Errorf("plaintext = %v, want %v", plaintext, tt.args.plaintext)
			}
		})
	}
}
