package http

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/k1nky/gophermart/internal/adapter/http/mock"
	"github.com/stretchr/testify/suite"

	"github.com/k1nky/gophermart/internal/entity/user"
)

type httpAdapterTestSuite struct {
	suite.Suite
	authService    *mock.MockauthService
	accountService *mock.MockaccountService
}

func TestHTTPAdapter(t *testing.T) {
	suite.Run(t, new(httpAdapterTestSuite))
}

func (suite *httpAdapterTestSuite) SetupTest() {
	ctrl := gomock.NewController(suite.T())
	suite.authService = mock.NewMockauthService(ctrl)
	suite.accountService = mock.NewMockaccountService(ctrl)
}

func (suite *httpAdapterTestSuite) TestRegister() {
	type want struct {
		statusCode          int
		authorizationHeader string
	}
	tests := []struct {
		name           string
		payload        string
		want           want
		expectRegister []interface{}
	}{
		{
			name:           "Valid",
			payload:        `{"login": "user", "password": "pass"}`,
			want:           want{statusCode: http.StatusOK, authorizationHeader: "sometoken"},
			expectRegister: []interface{}{"sometoken", nil},
		},
		{
			name:           "Invalid json",
			payload:        `{"login": "user", `,
			want:           want{statusCode: http.StatusBadRequest, authorizationHeader: ""},
			expectRegister: []interface{}{},
		},
		{
			name:           "Invalid body format",
			payload:        `{"login": "user", "pass": ""} `,
			want:           want{statusCode: http.StatusBadRequest, authorizationHeader: ""},
			expectRegister: []interface{}{},
		},
		{
			name:           "Duplicate login",
			payload:        `{"login": "user", "password": "somepassword"}`,
			want:           want{statusCode: http.StatusConflict, authorizationHeader: ""},
			expectRegister: []interface{}{"", user.ErrDuplicateLogin},
		},
		{
			name:           "Unexpected error",
			payload:        `{"login": "user", "password": "somepassword"}`,
			want:           want{statusCode: http.StatusInternalServerError, authorizationHeader: ""},
			expectRegister: []interface{}{"", errors.New("unexpected error")},
		},
	}
	a := &Adapter{
		auth: suite.authService,
	}
	for _, tt := range tests {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(tt.payload))
		if len(tt.expectRegister) > 0 {
			suite.authService.EXPECT().Register(gomock.Any(), gomock.Any()).Return(tt.expectRegister...)
		}
		a.Register(w, r)
		suite.Equal(tt.want.statusCode, w.Code)
	}
}

func (suite *httpAdapterTestSuite) TestLogin() {
	type want struct {
		statusCode          int
		authorizationHeader string
	}
	tests := []struct {
		name        string
		payload     string
		want        want
		expectLogin []interface{}
	}{
		{
			name:        "Valid",
			payload:     `{"login": "user", "password": "pass"}`,
			want:        want{statusCode: http.StatusOK, authorizationHeader: "sometoken"},
			expectLogin: []interface{}{"sometoken", nil},
		},
		{
			name:        "Invalid json",
			payload:     `{"login": "user", `,
			want:        want{statusCode: http.StatusBadRequest, authorizationHeader: ""},
			expectLogin: []interface{}{},
		},
		{
			name:        "Invalid body format",
			payload:     `{"login": "user", "pass": ""} `,
			want:        want{statusCode: http.StatusBadRequest, authorizationHeader: ""},
			expectLogin: []interface{}{},
		},
		{
			name:        "Invalid credentials",
			payload:     `{"login": "user", "password": "invalid_password"}`,
			want:        want{statusCode: http.StatusUnauthorized, authorizationHeader: ""},
			expectLogin: []interface{}{"", user.ErrInvalidCredentials},
		},
		{
			name:        "Unexpected error",
			payload:     `{"login": "user", "password": "somepassword"}`,
			want:        want{statusCode: http.StatusInternalServerError, authorizationHeader: ""},
			expectLogin: []interface{}{"", errors.New("unexpected error")},
		},
	}
	a := &Adapter{
		auth: suite.authService,
	}
	for _, tt := range tests {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(tt.payload))
		if len(tt.expectLogin) > 0 {
			suite.authService.EXPECT().Login(gomock.Any(), gomock.Any()).Return(tt.expectLogin...)
		}
		a.Login(w, r)
		suite.Equal(tt.want.statusCode, w.Code)
	}
}

func (suite *httpAdapterTestSuite) TestNewOrder() {
	type want struct {
		statusCode int
	}
	tests := []struct {
		name       string
		payload    string
		want       want
		mockExpect []interface{}
	}{
		{
			name:       "Success",
			payload:    "12345678903",
			want:       want{statusCode: http.StatusAccepted},
			mockExpect: []interface{}{nil, nil},
		},
		{
			name:       "Invalid number",
			payload:    "123456789",
			want:       want{statusCode: http.StatusUnprocessableEntity},
			mockExpect: []interface{}{},
		},
	}
	a := &Adapter{
		account: suite.accountService,
	}
	claims := user.PrivateClaims{
		ID:    user.ID(1),
		Login: "u1",
	}
	for _, tt := range tests {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(tt.payload))

		if len(tt.mockExpect) > 0 {
			suite.accountService.EXPECT().NewOrder(gomock.Any(), gomock.Any()).Return(tt.mockExpect...)
		}
		a.NewOrder(w, r.WithContext(context.WithValue(r.Context(), keyUserClaims, claims)))
		suite.Equal(tt.want.statusCode, w.Code)
	}
}
