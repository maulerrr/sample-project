package testing

import (
	"github.com/gin-gonic/gin"
	"github.com/maulerrr/sample-project/api/ctrl"
	"testing"
)

func TestAddLike(t *testing.T) {
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
			ctrl.AddLike(tt.args.context)
		})
	}
}

func TestGetLike(t *testing.T) {
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
			ctrl.GetLike(tt.args.context)
		})
	}
}

func TestGetLikesCountOnPost(t *testing.T) {
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
			ctrl.GetLikesCountOnPost(tt.args.context)
		})
	}
}
