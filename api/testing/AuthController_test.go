package testing

import (
	"github.com/gin-gonic/gin"
	"github.com/maulerrr/sample-project/api/ctrl"
	"github.com/maulerrr/sample-project/api/models"
	"testing"
)

func TestLogin(t *testing.T) {
	type args struct {
		context *gin.Context
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl.Login(tt.args.context)
		})
	}
}

func TestSignUp(t *testing.T) {
	type args struct {
		context *gin.Context
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl.SignUp(tt.args.context)
		})
	}
}

func Test_generateToken(t *testing.T) {
	type args struct {
		user models.User
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ctrl.generateToken(tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("generateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("generateToken() got = %v, want %v", got, tt.want)
			}
		})
	}
}
