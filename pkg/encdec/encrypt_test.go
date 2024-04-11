//go:build unit

package encdec

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
			name: "happy test",
			args: args{
				plaintext: "kingsman",
				key:       "04076d64bdb6fcf31706eea85ec98431"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// encrypt the plaintext
			ciphertext := Encrypt(tt.args.plaintext)

			t.Logf("ciphertext = %s", ciphertext)

			plaintext, err := Decrypt(ciphertext)
			if err != nil {
				t.Errorf("encrypt() error = %v", err)
				return
			}
			t.Logf("plaintext = %s", plaintext)
			//
			// compare the initial plaintext with output of previous decrypt function
			if plaintext != tt.args.plaintext {
				t.Errorf("plaintext = %v, want %v", plaintext, tt.args.plaintext)
			}
			//
		})
	}
}
