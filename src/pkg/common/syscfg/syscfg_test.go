package syscfg

import (
	"os"
	"testing"
    assert "github.com/deepinsight/mlcloud/src/utils/assert"
	"github.com/deepinsight/mlcloud/src/pkg/common"
)

func TestParseStringToInt(t *testing.T) {
	cases := []struct {
		input  string
		result int
	}{
		{"1", 1},
		{"-1", -1},
		{"0", 0},
		{"", 0},
	}

	for _, c := range cases {
		i, err := parseStringToInt(c.input)
		assert.Nil(t, err)
		assert.Equal(t, c.result, i)
	}
}

func TestParseStringToBool(t *testing.T) {
	cases := []struct {
		input  string
		result bool
	}{
		{"true", true},
		{"on", true},
		{"TRUE", true},
		{"ON", true},
		{"other", false},
		{"", false},
	}

	for _, c := range cases {
		b, _ := parseStringToBool(c.input)
		assert.Equal(t, c.result, b)
	}
}

func TestInitCfgStore(t *testing.T) {
	os.Clearenv()
	path := "/tmp/config.json"
	if err := os.Setenv("JSON_CFG_STORE_PATH", path); err != nil {
		t.Fatalf("failed to set env: %v", err)
	}
	defer os.RemoveAll(path)
	err := initCfgStore()
	assert.Nil(t, err)
}

func TestLoadFromEnv(t *testing.T) {
	os.Clearenv()
	ldapURL := "ldap://ldap.com"
	extEndpoint := "http://harbor.com"
	if err := os.Setenv("LDAP_URL", ldapURL); err != nil {
		t.Fatalf("failed to set env: %v", err)
	}
	cfgs := map[string]interface{}{}
	err := LoadFromEnv(cfgs, true)
	assert.Nil(t, err)
	assert.Equal(t, ldapURL, cfgs[common.LDAPURL])

	os.Clearenv()
	if err := os.Setenv("LDAP_URL", ldapURL); err != nil {
		t.Fatalf("failed to set env: %v", err)
	}
	if err := os.Setenv("EXT_ENDPOINT", extEndpoint); err != nil {
		t.Fatalf("failed to set env: %v", err)
	}

	cfgs = map[string]interface{}{}
	err = LoadFromEnv(cfgs, false)
	assert.Nil(t, err)
	assert.Equal(t, extEndpoint, cfgs[common.ExtEndpoint])
	assert.Equal(t, ldapURL, cfgs[common.LDAPURL])

	os.Clearenv()
	if err := os.Setenv("LDAP_URL", ldapURL); err != nil {
		t.Fatalf("failed to set env: %v", err)
	}
	if err := os.Setenv("EXT_ENDPOINT", extEndpoint); err != nil {
		t.Fatalf("failed to set env: %v", err)
	}

	cfgs = map[string]interface{}{
		common.LDAPURL: "ldap_url",
	}
	err = LoadFromEnv(cfgs, false)
	assert.Nil(t, err)
	assert.Equal(t, extEndpoint, cfgs[common.ExtEndpoint])
	assert.Equal(t, "ldap_url", cfgs[common.LDAPURL])
}
