/**
 * @Author: Hardews
 * @Date: 2023/3/9 17:51
 * @Description:
**/

package main

import (
	"spirit-core/api"
	"spirit-core/dao"
	"spirit-core/tool"
)

func main() {
	dao.InitDB()
	tool.RedisInit()
	tool.Cron()
	api.InitRouter()
}
