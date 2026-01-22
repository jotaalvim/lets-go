package models

import (
	"modulo.porreiro/internal/assert"
	"testing"
)

func TestUserModelExists(t *testing.T) {

	if testing.Short() {
		t.Skip("skipping long testss")
	}	



	tests := []struct {
		name   string
		userID int
		want   bool
	}{

		{
			name:   "ValidId",
			userID: 1,
			want:   true,
		},
		{
			name:   "Zero Id",
			userID: 0,
			want:   false,
		},
		{
			name:   "Non existent Id",
			userID: 2,
			want:   false,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			db := newTestDB(t)

			m := UserModel{db}

			exists, err := m.Exists(tt.userID)
			assert.Equal(t, exists, tt.want)
			assert.NilError(t, err)
		})
	}
}
