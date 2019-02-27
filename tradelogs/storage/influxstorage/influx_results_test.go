package influxstorage

import (
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"
)

func loadTestData(db string) error {
	const endpoint = "http://127.0.0.1:8086/"
	client := http.Client{Timeout: time.Second * 5}

	req, err := http.NewRequest(http.MethodPost, endpoint+"query", nil)
	if err != nil {
		return err
	}

	q := req.URL.Query()
	q.Add("q", fmt.Sprintf("CREATE DATABASE %s", db))
	req.URL.RawQuery = q.Encode()

	rsp, err := client.Do(req)
	if err != nil {
		return err
	}
	if rsp.StatusCode != http.StatusOK {
		return fmt.Errorf("wrong status code, expected: %d, got: %d", http.StatusOK, rsp.StatusCode)
	}

	exported, err := os.Open("../testdata/export.dat")
	if err != nil {
		return err
	}
	defer exported.Close()

	req, err = http.NewRequest(http.MethodPost, endpoint+"write", exported)
	if err != nil {
		return err
	}

	q = req.URL.Query()
	q.Add("db", db)
	req.URL.RawQuery = q.Encode()
	req.Header.Set("Content-Type", "binary/octet-stream")

	rsp, err = client.Do(req)
	if err != nil {
		return err
	}

	if rsp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("wrong status code, expected: %d, got: %d", http.StatusNoContent, rsp.StatusCode)
	}
	return nil
}

func TestLoadTradeLogs(t *testing.T) {
	const (
		dbName = "test_results"
		// These params must be change when export.dat changes.
		fromTime    = 1539247511
		toTime      = 1539248681
		expectedLen = 11
	)

	is, err := newTestInfluxStorage(dbName)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err = is.tearDown(); err != nil {
			t.Fatal(err)
		}
	}()
	if err = loadTestData(dbName); err != nil {
		t.Fatal(err)
	}

	// following verification can be changed once the export.dat file is regen.
	tradeLogs, err := is.LoadTradeLogs(time.Unix(fromTime, 0), time.Unix(toTime, 0))
	if err != nil {
		t.Fatal(err)
	}

	if len(tradeLogs) != expectedLen {
		t.Errorf("wrong number of trade log returned, expected: %d, got: %d", expectedLen, len(tradeLogs))
	}

}
