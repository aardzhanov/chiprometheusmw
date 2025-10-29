package prometheusmw

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/stretchr/testify/require"
)

func TestMiddleware(t *testing.T) {

	testCases := []struct {
		name string
		path []string
		want string
	}{
		{
			name: "Test SimplePath",
			path: []string{"simplepath"},
			want: `code="200",method="GET",path="/simplepath",service="test"`,
		},
		{
			name: "Tst Path with macros",
			path: []string{"pathwithmacro", "value1"},
			want: `code="200",method="GET",path="/pathwithmacro/{macro}",service="test"`,
		},
	}

	recorder := httptest.NewRecorder()
	router := chi.NewRouter()
	metrics := NewMiddleware("test")
	router.Use(metrics.Handler)
	router.Get(`/simplepath`, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	router.Get(`/pathwithmacro/{macro}`, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	router.Handle("/metrics", promhttp.Handler())

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			request, err := http.NewRequest("GET", "http://localhost/"+strings.Join(testCase.path, "/"), nil)
			require.NoError(t, err)
			router.ServeHTTP(recorder, request)
			metric, err := http.NewRequest("GET", "http://localhost/metrics", nil)
			require.NoError(t, err)
			router.ServeHTTP(recorder, metric)
			require.Contains(t, recorder.Body.String(), testCase.want)
		})

	}
}
