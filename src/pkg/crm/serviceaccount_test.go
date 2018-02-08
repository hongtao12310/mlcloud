package crm

import (
    "github.com/deepinsight/mlcloud/src/pkg/models"
    "testing"
)

func TestCreateServiceAccount(t *testing.T)  {
    user := &models.User{
        Username: "test-user",
        Sysadmin: false,
    }

    // create namespace
    if err := CreateUserNamespace(user); err != nil {
        t.Fatal(err)
    }

    // create rolebinding
    if err := CreateUserRolebinding(user); err != nil {
        t.Fatal(err)
    }

    // create mlcloud service account
    if err := CreateMLCloudServiceAccount(user); err != nil {
        t.Fatal(err)
    }

    // list pods
    if err := ListPod(user); err != nil {
        t.Fatal(err)
    }
}
