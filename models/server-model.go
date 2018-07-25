package models

type Server struct {
	Id         int64  `json:"id"`
	Ip         string `json:"ip"`
	CoreNum    string `json:"coreNum" form:"coreNum"`
	MemorySize string `json:"memorySize" form:"memorySize"`
	CreateAt   Time
	UpdateAt   Time
}

type ServerDetail struct {
	Server
	Port       int
	Username   string
	Password   string
	PrivateKey string
}

func ListServer() []Server {
	results := []Server{}
	rows, err := Mgr.db.Query("SELECT id,ip,core_num,memory_size FROM server")
	defer rows.Close()
	checkErr(err)

	for rows.Next() {
		var server Server
		err = rows.Scan(
			&server.Id, &server.Ip, &server.CoreNum, &server.MemorySize)
		checkErr(err)
		results = append(results, server)
	}
	return results
}

func GetServerDetail(ip string) ServerDetail {
	var d ServerDetail
	sql := "SELECT id,ip,core_num,memory_size,username,password,private_key,port,created_at,updated_at FROM server WHERE ip = ?"
	row := Mgr.db.QueryRow(sql, ip)
	err := row.Scan(&d.Id, &d.Ip, &d.CoreNum, &d.MemorySize, &d.Username, &d.Password, &d.PrivateKey, &d.Port, &d.CreateAt, &d.UpdateAt)
	checkErr(err)
	return d
}
