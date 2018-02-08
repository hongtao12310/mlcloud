package json

import (
	"os"
	"testing"
    "reflect"
)

func TestReadWrite(t *testing.T) {
	path := "/tmp/config.json"
	store, err := NewCfgStore(path)
	if err != nil {
		t.Fatalf("failed to create json cfg store: %v", err)
	}
	defer func() {
		if err := os.Remove(path); err != nil {
			t.Fatalf("failed to remove the json file %s: %v", path, err)
		}
	}()

	if store.Name() != "JSON" {
		t.Errorf("unexpected name: %s != %s", store.Name(), "JSON")
		return
	}

	config := map[string]interface{}{
		"mysql_host": "localhost",
        "mysql_port": 3306,
	}

	if err := store.Write(config); err != nil {
		t.Errorf("failed to write configurations to json file: %v", err)
		return
	}

    output, err := store.Read()

	if err != nil {
		t.Errorf("failed to read configurations from json file: %v", err)
		return
	}

    t.Logf("mysql_host type: %s", reflect.TypeOf(output["mysql_host"]))
    t.Logf("mysql_port type: %s", reflect.TypeOf(output["mysql_port"]))
    t.Logf("config output: %+v", output)
}
