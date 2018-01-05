package httpClient

import (
	"testing"
)

func TestClient(t *testing.T) {
	c := NewClient()

	c.SetTimeout(10)
	c.SetHeader("Tester", "OK")

	h := map[string]string{
		"Content-Type": "application/json",
		"Channel":      "web",
	}

	c.SetHeaders(h)

	if c.timeout != 10 {
		t.Errorf("Expected 10, got : %d", c.timeout)
	}

	if c.headers["Tester"] != "OK" {
		t.Errorf("Expected OK, got : %s", c.headers["Tester"])
	}

	for x, y := range h {

		if c.headers[x] != y {
			t.Errorf("Expected %s, got : %s", y, c.headers[x])
		}
	}



}
