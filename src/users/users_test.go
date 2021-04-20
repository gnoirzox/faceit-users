package users

import "testing"

func TestIsValidNickname(t *testing.T) {
	var tests = []struct {
		input string
		want  bool
	}{
		{"sim", true},
		{"simulation123456", false},
		{"lo", false},
		{"ye", false},
		{"l", false},
		{"bulldog", true},
		{"dog", true},
		{"catfish", true},
		{"CATfish", true},
		{"catFISH", true},
		{"seahorse", true},
		{"seaHORSE", true},
		{"SEAhorse", true},
	}

	for _, test := range tests {
		var u User
		u.Nickname = test.input

		if got := u.IsValidNickname(); got != test.want {
			t.Errorf("u.Username = %q u.IsValidNickname() = %v", test.input, got)
		}
	}
}

func TestIsValidPassword(t *testing.T) {
	var tests = []struct {
		input string
		want  bool
	}{
		{"yo", false},
		{"heifei", false},
		{"MyPaSS0899!W)rd", true},
		{"MyLongPassword14327854", true},
	}

	for _, test := range tests {
		var u User
		u.Password = test.input

		if got := u.IsValidPassword(); got != test.want {
			t.Errorf("u.Password = %q u.IsValidPassword() = %v", test.input, got)
		}
	}
}

func TestIsValidEmail(t *testing.T) {
	var tests = []struct {
		input string
		want  bool
	}{
		{"x@test.com", true},
		{"simon.rouger@test.com", true},
		{"sim..@test.com", false},
		{"sim@@test.com", false},
		{"sim@test@com", false},
		{"sim@test..com", false},
	}

	for _, test := range tests {
		var u User
		u.Email = test.input

		if got := u.IsValidEmail(); got != test.want {
			t.Errorf("u.Email = %q u.IsValidEmail() = %v", test.input, got)
		}
	}
}
