package models

import (
	"fmt"
)

type UserInfo struct {
	Id        int64
	Username  string
	Email     string
	Ticket    string
	Role      string
	Enable    string
	CreatedAt Time
	UpdatedAt Time
}

func DetailByEmail(email string) (UserInfo, error) {
	sql := "SELECT id,username,email,ticket,role,enable,created_at,updated_at FROM userinfo where email = ?"
	row := Mgr.db.QueryRow(sql, email)
	var u UserInfo
	if err := row.Scan(&u.Id, &u.Username, &u.Email, &u.Ticket, &u.Role, &u.Enable, &u.CreatedAt, &u.UpdatedAt); err != nil {
		fmt.Errorf("%s find err %v", email, err)
		return u, err
	}
	return u, nil
}
