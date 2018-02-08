package fs

import (
    "github.com/deepinsight/mlcloud/src/pkg/models"
    "github.com/deepinsight/mlcloud/src/pkg/common/syscfg"
    "github.com/deepinsight/mlcloud/src/pkg/common"
    "os"
    "path"
    "github.com/golang/glog"
)

// create the user base directory for user in the file system
func CreateBaseDirForUser(user *models.User) error {
    config := syscfg.GetSysConfig()
    fsBaseDir := config.Get(common.FSBasePath).(string)
    userDir := path.Join(fsBaseDir, user.Username)
    if _, err := os.Stat(userDir); os.IsNotExist(err) {
        if err := os.MkdirAll(userDir, 0755); err != nil {
            return err
        }
    }
    return nil
}

// remove user base dir from the file system
func DeleteBaseDirForUser(user *models.User) error {
    glog.V(4).Infof("DeleteBaseDirForUser")
    config := syscfg.GetSysConfig()
    fsBaseDir := config.Get(common.FSBasePath).(string)
    userDir := path.Join(fsBaseDir, user.Username)

    _, err := os.Stat(userDir)
    //glog.V(4).Infof("err: %v", err)
    //glog.V(4).Infof("os.IsExist(err): %v", os.IsExist(err))
    if err == nil {
        glog.V(4).Infof("%s exists", userDir)
        if err := os.RemoveAll(userDir); err != nil {
            return err
        }
    }

    glog.V(4).Infof("DeleteBaseDirForUser: end")
    return nil
}