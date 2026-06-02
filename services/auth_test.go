package services

import (
	"os"
	"testing"

	"expense-tracker-api/models"

	beego "github.com/beego/beego/v2/server/web"
)

func TestMain(m *testing.M) {
	// move working directory to project root so beego finds conf/app.conf
	os.Chdir("../")

	// point to test-specific files so real data is never touched
	beego.AppConfig.Set("datadir", "testdata")
	beego.AppConfig.Set("userfile", "users_test.csv")
	beego.AppConfig.Set("expensefile", "expenses_test.csv")

	code := m.Run()

	// clean up after all tests finish
	os.RemoveAll("testdata")

	os.Exit(code)
}

// cleanTestData removes test CSV files before each test run
func cleanTestData() {
	os.RemoveAll("testdata")
	os.MkdirAll("testdata", os.ModePerm)
}

func TestRegisterUser(t *testing.T) {
	test := []struct {
		name    string
		input models.RegisterRequest
		wantErr bool
		errMsg string
	} {
		{
			name: "valid registration",
			input: models.RegisterRequest{
				Name: "John Doe",
				Email: "john@example.com",
				Password: "secret123",
			},
			wantErr: false,
		},
		{
			name: "duplicate email",
			input: models.RegisterRequest{
				Name: "Jane Doe",
				Email: "john@example.com",
				Password: "anotherpassword",
			},
			wantErr: true,
			errMsg: "email already exists",
		}, 
		{
			name: "different email succeeds",
			input: models.RegisterRequest{
				Name: "Jane Doe",
				Email: "jane@example.com",
				Password: "secret123",
			},
			wantErr: false,
		},
	}

	cleanTestData()

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			err := RegisterUser(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error %q, got nil", tt.errMsg)
					return
				}
				if err.Error() !=  tt.errMsg {
					t.Errorf("expected error %q, got %q", tt.errMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", err.Error())
				}
			}
		})
	}
}

func TestLoginUser(t *testing.T) {
	// seed one user before login tests
	cleanTestData()
	RegisterUser(models.RegisterRequest{
		Name: "John Doe",
		Email: "john@example.com",
		Password: "secret123",
	})

	tests := []struct {
		name string
		input models.LoginRequest
		wantErr bool
		wantID bool
	}{
		{
			name: "valid credentials",
			input: models.LoginRequest{
				Email: "john@example.com",
				Password: "secret123",
			},
			wantErr: false,
			wantID: true,
		},
		{
			name: "wrong password",
			input: models.LoginRequest{
				Email: "john@example.com",
				Password: "wrongpassword",
			},
			wantErr: true,
		},
		{
			name: "email not found",
			input: models.LoginRequest{
				Email: "nonexistent@example.com",
				Password: "secret123",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := LoginUser(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("expected no error, got %v", err.Error())
				return
			}
			if tt.wantID && data.UserID == 0 {
				t.Errorf("expected valid user ID, got 0")
			}
		})
	}
}