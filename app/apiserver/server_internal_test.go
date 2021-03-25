package apiserver

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/iavealokin/cashflow/app/store/teststore"
	"github.com/stretchr/testify/assert"
)

//TestServerHandleUsersCreate ...
func TestServerHandleUsersCreate(t *testing.T) {
	s := newServer(teststore.New())
	testCases:=[]struct{
		name string
		payload string
		expectedCode int
	}{
		{
			name: "valid",
			payload: `{
				"login":"evgen",
  				"username":"Evgeniy",
  				"surname":"Sheblyakin",
  				"birthday":"13.09.1794",
 				"password":"SanyaIsTheBest"
			}`,
			expectedCode: http.StatusCreated,
		},
		{
			name: "invalid payload",
			payload: "invalid",
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "mismathced types",
			payload:`{
				"login":5,
  				"username":"Evgeniy",
  				"surname":"Sheblyakin",
  				"birthday":"13.09.1794",
				"password":"SanyaIsTheBest",
				"userid": "string"
			}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "Lentgh less 3 symbols",
			payload:`{
				"login":"2",
  				"username":"Evgeniy",
  				"surname":"Sheblyakin",
  				"birthday":"13.09.1794",
				"password":"SanyaIsTheBest"
			}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "Empty fields",
			payload:`{
				"login":"",
  				"username":"Evgeniy",
  				"surname":"Sheblyakin",
  				"birthday":"13.09.1794",
				"password":"SanyaIsTheBest"
			}`,
			expectedCode: http.StatusBadRequest,
		},
	}


	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T){
			rec := httptest.NewRecorder()
		b:= strings.NewReader(tc.payload)
			req, _ := http.NewRequest(http.MethodPost, "/addUser",b)
			s.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})

	}
}


//TestServerHandleUsersUpdate ...
func TestServerHandleUsersUpdate(t *testing.T) {
	s := newServer(teststore.New())
	testCases:=[]struct{
		name string
		payload string
		expectedCode int
	}{
		{
			name: "valid",
			payload: `{
				"login":"sanek",
				"username":"Alexandr",
				"surname":"Penykov",
				"birthday":"13.09.1994",
				"password":"SanyaIsTheBest",
				"userid":1
			}`,
			expectedCode: http.StatusCreated,
		},
		{
			name: "invalid payload",
			payload: "invalid",
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "mismathced types",
			payload:`{
				"login":5,
  				"username":"Evgeniy",
  				"surname":"Sheblyakin",
  				"birthday":"13.09.1794",
				"password":"SanyaIsTheBest",
				"userid": "string"
			}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "Lentgh less 3 symbols",
			payload:`{
				"login":"2",
  				"username":"Evgeniy",
  				"surname":"Sheblyakin",
  				"birthday":"13.09.1794",
				"password":"SanyaIsTheBest"
			}`,
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "Empty fields",
			payload:`{
				"login":"",
  				"username":"Evgeniy",
  				"surname":"Sheblyakin",
  				"birthday":"13.09.1794",
				"password":"SanyaIsTheBest"
			}`,
			expectedCode: http.StatusBadRequest,
		},
	}


	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T){
			rec := httptest.NewRecorder()
		b:= strings.NewReader(tc.payload)
			req, _ := http.NewRequest(http.MethodPost, "/updateUser",b)
			s.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})

	}
}

