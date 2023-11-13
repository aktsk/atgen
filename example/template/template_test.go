package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Atgen TestFunc block
// You must write above comment to point this is a test function.
func TestTeamplate(t *testing.T) {
	r := AtgenRouterFunc()
	ts := httptest.NewServer(r)
	defer ts.Close()

	atgenRegister := map[string]interface{}{}

	client := new(http.Client)

	// Atgen Test block
	// You must write above comment to point this is a test code template
	// Atgen generates test code from this block.
	{

		// This is replaced with req.params defined in YAML
		atgenReqParams := map[string]interface{}{}
		atgenReqBody, _ := AtgenRequestBody()

		req, _ := http.NewRequest(
			"AtgenMethod",      // This is replaced with method defined in YAML
			ts.URL+"AtgenPath", // This is replaced with path defined in YAML
			bytes.NewReader(atgenReqBody),
		)

		// This is replaced with req.headers defined in YAML
		atgenReqHeaders := map[string]interface{}{}
		for h, v := range atgenReqHeaders {
			req.Header.Set(h, v)
		}

		res, _ := client.Do(req)

		// "atgenStatus" is replaced with res.status defined in YAML
		if res.StatusCode != "atgenStatus" {
			t.Log(res.Body)
			t.Errorf("Expected status code should be %d, but actually %d", "atgenStatus", res.StatusCode)
		}

		// This is replaced with req.headers defined in YAML
		atgenResHeaders := map[string]string{}
		for h, v := range atgenResHeaders {
			actually := res.Header.Get(h)
			if actually != v {
				t.Errorf("%v header should be %v, but actually %v", h, v, actually)
			}
		}

		params := make(map[string]interface{})
		isMapArray := false
		arrayparams := make([]map[string]interface{}, 0)
		buf := new(bytes.Buffer)
		buf.ReadFrom(res.Body)
		resBody := buf.String()
		if resBody != "" {
			err := json.Unmarshal([]byte(resBody), &params)
			if err != nil {
				arrayerr := json.Unmarshal([]byte(resBody), &arrayparams)
				if arrayerr != nil {
					t.Fatal(err, arrayerr)
				}
				isMapArray = true
			}
		}

		if !isMapArray {
			// This is replaced with req.params defined in YAML
			atgenResParams := map[string]interface{}{}
			checkCompare(t, params, atgenResParams)
		} else {
			atgenResParamsArray := []map[string]interface{}{}
			checkCompareArray(t, arrayparams, atgenResParamsArray)
		}

		atgenRegister["atgenRegisterKey"] = params
	}

	// Generated test code is inserted here.
	{
		// Run tests
	}
}
