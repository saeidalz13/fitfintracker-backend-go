package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	cn "github.com/saeidalz13/LifeStyle2/lifeStyleBack/config"
	sqlc "github.com/saeidalz13/LifeStyle2/lifeStyleBack/db/sqlc"
)

func TestSignUp(t *testing.T) {
	app.Post(cn.URLS.SignUp, testHandlersManager.AuthHandler.HandlePostSignUp)

	tests := []TestCase[sqlc.CreateUserParams]{
		{
			name:           "signup valid email and password",
			url:            cn.URLS.SignUp,
			expectedStatus: http.StatusCreated,
			contentType:    validContentType,
			reqPayload: sqlc.CreateUserParams{
				Email:    validEmail,
				Password: validPassword,
			},
		},
		{
			name:           "invalid email",
			url:            cn.URLS.SignUp,
			expectedStatus: http.StatusConflict,
			contentType:    validContentType,
			reqPayload: sqlc.CreateUserParams{
				Email:    "invalid",
				Password: validPassword,
			},
		},
		{
			name:           "not signup with existing email",
			url:            cn.URLS.SignUp,
			expectedStatus: http.StatusConflict,
			contentType:    validContentType,
			reqPayload: sqlc.CreateUserParams{
				Email:    validEmail,
				Password: validPassword,
			},
		},
		{
			name:           "invalid short password",
			url:            cn.URLS.SignUp,
			expectedStatus: http.StatusConflict,
			contentType:    validContentType,
			reqPayload: sqlc.CreateUserParams{
				Email:    anotherValidEmail,
				Password: "short",
			},
		},
	}

	for i, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			body, err := json.Marshal(test.reqPayload)
			if err != nil {
				t.Fatal(err)
			}

			req := httptest.NewRequest(http.MethodPost, test.url, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", test.contentType)

			resp, err := app.Test(req, 10)
			if err != nil {
				t.Fatal(err)
			}

			if resp.StatusCode != test.expectedStatus {
				t.Fatalf("expected status code: %d\t got: %d", resp.StatusCode, test.expectedStatus)
			}

			if i == 0 {
				testMock.ExpectQuery(`SELECT COUNT(id) from users where email = \$1`).
				WithArgs(validEmail).
				WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

				if err := testMock.ExpectationsWereMet(); err != nil {
					t.Fatal(err)
				}
			}
		})

	}
}

func TestLogin(t *testing.T) {
	app.Post(cn.URLS.Login, testHandlersManager.AuthHandler.HandlePostLogin)

	tests := []TestCase[sqlc.CreateUserParams]{
		{
			name:           "login with valid email and password",
			url:            cn.URLS.Login,
			expectedStatus: http.StatusOK,
			contentType:    validContentType,
			reqPayload: sqlc.CreateUserParams{
				Email:    validEmail,
				Password: validPassword,
			},
		},
		{
			name:           "no login with invalid password",
			url:            cn.URLS.Login,
			expectedStatus: http.StatusUnauthorized,
			contentType:    validContentType,
			reqPayload: sqlc.CreateUserParams{
				Email:    validEmail,
				Password: "invalid",
			},
		},
		{
			name:           "no login with invalid email",
			url:            cn.URLS.Login,
			expectedStatus: http.StatusUnauthorized,
			contentType:    validContentType,
			reqPayload: sqlc.CreateUserParams{
				Email:    "invalid",
				Password: validPassword,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			body, err := json.Marshal(test.reqPayload)
			if err != nil {
				t.Fatal(err)
			}

			req := httptest.NewRequest(http.MethodPost, test.url, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", test.contentType)
			resp, err := app.Test(req, 10)
			if err != nil {
				t.Fatal(err)
			}

			if resp.StatusCode != test.expectedStatus {
				t.Fatalf("expected status code: %d\t got: %d", resp.StatusCode, test.expectedStatus)
			}
		})
	}
}

func TestGetHome(t *testing.T) {
	app.Get(cn.URLS.Home, testHandlersManager.AuthHandler.HandleGetHome)

	tests := []TestCase[sqlc.CreateUserParams]{
		{
			name:           "get home valid token",
			url:            cn.URLS.Home,
			expectedStatus: http.StatusOK,
			token:          validPasetoToken,
		},
		{
			name:           "get home invalid token unauthorized",
			url:            cn.URLS.Home,
			expectedStatus: http.StatusUnauthorized,
			token:          invalidToken,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, test.url, nil)
			req.AddCookie(&http.Cookie{Name: "paseto", Value: test.token})

			resp, err := app.Test(req, 10)
			if err != nil {
				t.Fatal(err)
			}

			if resp.StatusCode != test.expectedStatus {
				t.Fatalf("expected status code: %d\t got: %d", resp.StatusCode, test.expectedStatus)
			}
		})
	}
}

func TestGetProfile(t *testing.T) {
	app.Get(cn.URLS.Profile, testHandlersManager.AuthHandler.HandleGetProfile)

	tests := []TestCase[sqlc.CreateUserParams]{
		{
			name:           "get profile valid token",
			url:            cn.URLS.Profile,
			expectedStatus: http.StatusOK,
			token:          validPasetoToken,
		},
		{
			name:           "no profile invalid token unauthorized",
			url:            cn.URLS.Profile,
			expectedStatus: http.StatusUnauthorized,
			token:          invalidToken,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, test.url, nil)
			req.AddCookie(&http.Cookie{Name: "paseto", Value: test.token})

			resp, err := app.Test(req, 10)
			if err != nil {
				t.Fatal(err)
			}

			if resp.StatusCode != test.expectedStatus {
				t.Fatalf("expected status code: %d\t got: %d", resp.StatusCode, test.expectedStatus)
			}
		})
	}
}

func TestDeleteUser(t *testing.T) {
	tests := []TestCase[sqlc.CreateUserParams]{
		{
			name:           "should not delete user invalid token",
			url:            cn.URLS.DeleteProfile,
			expectedStatus: http.StatusUnauthorized,
			token:          invalidToken,
		},
		{
			name:           "should delete user valid token",
			url:            cn.URLS.DeleteProfile,
			expectedStatus: http.StatusNoContent,
			token:          validPasetoToken,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodDelete, test.url, nil)
			req.AddCookie(&http.Cookie{Name: "paseto", Value: test.token})

			resp, err := app.Test(req, 10)
			if err != nil {
				t.Fatal(err)
			}

			if resp.StatusCode != test.expectedStatus {
				t.Fatalf("expected status code: %d\t got: %d", resp.StatusCode, test.expectedStatus)
			}

			if test.name == "should delete user valid token" {
				testMock.ExpectQuery(`SELECT COUNT(*) from users where email = \$1`).
				WithArgs(validEmail).
				WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

				if err := testMock.ExpectationsWereMet(); err != nil {
					t.Fatal(err)
				}
			}
		})
	}
}
