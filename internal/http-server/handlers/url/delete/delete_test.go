package delete_test

import (
	"UrlShort/internal/http-server/handlers/slogdiscard"
	"UrlShort/internal/http-server/handlers/url/delete"
	"UrlShort/internal/http-server/handlers/url/delete/mocks"
	"UrlShort/internal/lib/api/response"
	"UrlShort/internal/storage"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeleteHandler(t *testing.T) {
	cases := []struct {
		name         string
		alias        string
		mockError    error
		respError    string
		expectStatus int
	}{
		{
			name:         "Success",
			alias:        "test_alias",
			expectStatus: http.StatusOK,
		},
		{
			name:         "Alias Not Found",
			alias:        "nonexistent_alias",
			mockError:    storage.ErrURLNotFound,
			respError:    "not found",
			expectStatus: http.StatusNotFound,
		},
		{
			name:         "Internal Error",
			alias:        "test_alias",
			mockError:    errors.New("internal error"),
			respError:    "internal error",
			expectStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			urlDeleterMock := mocks.NewURLDeleter(t)

			if tc.alias != "" {
				if tc.mockError != nil {
					urlDeleterMock.On("DeleteURL", tc.alias).Return(tc.mockError).Once()
				} else {
					urlDeleterMock.On("DeleteURL", tc.alias).Return(nil).Once()
				}
			}

			r := chi.NewRouter()
			r.Delete("/{alias}", delete.New(slogdiscard.NewDiscardLogger(), urlDeleterMock))

			ts := httptest.NewServer(r)
			defer ts.Close()

			req, err := http.NewRequest(http.MethodDelete, ts.URL+"/"+tc.alias, nil)
			require.NoError(t, err)

			resp, err := http.DefaultClient.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, tc.expectStatus, resp.StatusCode)

			if resp.Header.Get("Content-Type") == "application/json" {
				var apiResponse response.Response
				require.NoError(t, json.NewDecoder(resp.Body).Decode(&apiResponse))
				assert.Equal(t, tc.respError, apiResponse.Error)
			}

			if tc.alias == "" {
				urlDeleterMock.AssertNotCalled(t, "DeleteURL")
			} else {
				urlDeleterMock.AssertExpectations(t)
			}
		})
	}
}
