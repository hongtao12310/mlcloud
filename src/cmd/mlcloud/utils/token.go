package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
    "github.com/deepinsight/mlcloud/src/utils/log"
    "github.com/deepinsight/mlcloud/src/utils"
    "github.com/deepinsight/mlcloud/src/utils/http"
)

// login MLCloud and get the user token cache
func getToken(uri string, body []byte) ([]byte, error) {
	req, err := MakeRequest(uri, "POST", bytes.NewBuffer(body), "", nil, nil)
	if err != nil {
		return nil, err
	}
	return GetResponse(req)
}

// Token fetch and caches the token for current configured user
func Token(config *SubmitConfig) (string, error) {
	tokenPath := filepath.Join(utils.UserHomeDir(), ".mlcloud", "token_cache")
	tokenbytes, err := ioutil.ReadFile(tokenPath)
	if err != nil {
		log.Warning("previous token not found, fetching a new one...")
		// Authenticate to the cloud endpoint
		authJSON := map[string]string{}
		authJSON["username"] = config.ActiveConfig.Username
		authJSON["password"] = config.ActiveConfig.Password
		authStr, _ := json.Marshal(authJSON)

		tokenURI := config.ActiveConfig.Endpoint + "/api/v1/login"
		log.Infof("token URI: %s", tokenURI)
		body, err := getToken(tokenURI, authStr)
		if err != nil {
			return "", err
		}

		respObj := http.Response{}
		if errJSON := json.Unmarshal(body, &respObj); errJSON != nil {
			return "", errJSON
		}

		if respObj.Err != "" {
			log.Fatal(respObj.Err)
		}

		log.Infof("response object: %+v", respObj)

		tokenStr := respObj.Results.(map[string]interface{})["token"].(string)
		err = ioutil.WriteFile(tokenPath, []byte(tokenStr), 0600)
		if err != nil {
			fmt.Fprintf(os.Stderr, "write cache token file error: %v", err)
		}
		// Ignore write token error, fetch a new one next time
		return tokenStr, nil
	}
	return string(tokenbytes), nil
}
