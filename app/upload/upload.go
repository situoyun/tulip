package upload

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"git.autops.xyz/autops/base/database/sql"
	"git.autops.xyz/autops/base/logs"
	"git.autops.xyz/autops/tulip/app"
)

var path = "/opt/"

func Run() {
	jsonFile, err := os.Open(path + "user.json")
	if err != nil {
		logs.Errorf("run| open file failed:%v", err)
		return
	}
	defer jsonFile.Close()
	fileData, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		logs.Errorf("run| read file failed:%v", err)
		return
	}

	userM := make(map[uint]*app.User)
	err = json.Unmarshal(fileData, &userM)
	if err != nil {
		logs.Errorf("run| json.unmarsh  failed:%v", err)
		return
	}

	for _, u := range userM {
		param := make(map[string]interface{})
		param["password"] = u.PassWord
		param["password_type"] = u.PassWordType
		param["salt"] = u.Salt

		if err := sql.GetDB().Model(app.User{}).Where("id = ?", u.ID).Updates(param).Error; err != nil {
			logs.Errorf("run|update user[%d] failed:%v", u.ID, err)
			continue
		}
	}
}
