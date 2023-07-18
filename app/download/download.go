package download

import (
	"encoding/json"
	"os"

	"git.autops.xyz/autops/base/database/sql"
	"git.autops.xyz/autops/base/logs"
	"git.autops.xyz/autops/tulip/app"
)

var path = "/opt/"

func Run() {
	var users []*app.User
	if err := sql.GetDB().Model(app.User{}).Scan(&users).Error; err != nil {
		logs.Errorf("run| get user failed:%v", err.Error())
		return
	}

	userM := make(map[uint]*app.User)
	for _, u := range users {
		userM[u.ID] = u
	}

	data, _ := json.Marshal(userM)

	err := os.WriteFile(path+"user.json", data, 0644)
	if err != nil {
		logs.Errorf("run| open file failed:%v", err.Error())
		return
	}
}
