package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestFunc block
// You must write above comment to point this is a test function.
func TestTeamplate(t *testing.T) {
	r := RouterFunc()
	ts := httptest.NewServer(r)
	defer ts.Close()

	client := new(http.Client)

	// Test block
	// You must write above comment to point this is a test code template
	// Atgen generates test code from this block.
	{

		// This is replaced with req.params defined in YAML
		reqParams := map[string]interface{}{}
		reqBody, _ := json.Marshal(reqParams)

		req, _ := http.NewRequest(
			"Method",      // This is replaced with method defined in YAML
			ts.URL+"Path", // This is replaced with path defined in YAML
			bytes.NewReader(reqBody),
		)

		// This is replaced with req.headers defined in YAML
		reqHeaders := map[string]interface{}{}
		for h, v := range reqHeaders {
			req.Header.Set(h, v)
		}

		res, _ := client.Do(req)

		// "status" is replaced with res.status defined in YAML
		if res.StatusCode != "status" {
			t.Log(res.Body)
			t.Errorf("Expected status code should be %d, but actually %d", "status", res.StatusCode)
		}

		// This is replaced with req.headers defined in YAML
		resHeaders := map[string]string{}
		for h, v := range resHeaders {
			actually := res.Header.Get(h)
			if actually != v {
				t.Errorf("%v header should be %v, but actually %v", h, v, actually)
			}
		}

		params := make(map[string]interface{})

		buf := new(bytes.Buffer)
		buf.ReadFrom(res.Body)
		resBody := buf.String()
		if resBody != "" {
			err := json.Unmarshal([]byte(resBody), &params)
			if err != nil {
				t.Fatal(err)
			}
		}

		// This is replaced with req.params defined in YAML
		resParams := map[string]interface{}{}
		for k, v := range resParams {
			if params[k] != v {
				t.Fatalf("params[%#v] should be %#v, but %#v", k, v, params[k])
			}
		}
	}

	// Generated test code is inserted here.
	{
		// Run tests
	}
}
