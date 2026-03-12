package crypto

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncrypt(t *testing.T) {
	type args struct {
		plainPassword string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Encrypt",
			args: args{
				plainPassword: "112233",
			},
			want:    "$2a$10$NoZtpzNfskZdbZSUzWWTV.slfr4/cKQsg5SNrzlxOHehPTn60.l4m",
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Encrypt(tt.args.plainPassword)
			fmt.Print(got)
			if !tt.wantErr(t, err, fmt.Sprintf("Encrypt(%v)", tt.args.plainPassword)) {
				return
			}
			assert.Equalf(t, tt.want, got, "Encrypt(%v)", tt.args.plainPassword)
		})
	}
}

func TestVerify(t *testing.T) {
	type args struct {
		plainPassword string
		encrypted     string
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Verify",
			args: args{
				plainPassword: "112233",
				encrypted:     "$2a$10$NoZtpzNfskZdbZSUzWWTV.slfr4/cKQsg5SNrzlxOHehPTn60.l4m",
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.wantErr(t, Verify(tt.args.plainPassword, tt.args.encrypted), fmt.Sprintf("Verify(%v, %v)", tt.args.plainPassword, tt.args.encrypted))
		})
	}
}
