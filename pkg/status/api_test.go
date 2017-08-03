package status

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAPI(t *testing.T) {
	t.Parallel()

	t.Run("status", func(t *testing.T) {
		var (
			api    = NewAPI()
			server = httptest.NewServer(api)
		)
		defer server.Close()

		response, err := http.Get(server.URL)
		if err != nil {
			t.Fatal(err)
		}

		if expected, actual := http.StatusOK, response.StatusCode; expected != actual {
			t.Errorf("expected: %d, actual: %d", expected, actual)
		}
	})
}
