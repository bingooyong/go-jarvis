package models

import (
	_ "github.com/mattn/go-sqlite3"
	"testing"
	"github.com/sirupsen/logrus"
)

func TestGetServerDetail(t *testing.T) {
	var detail ServerDetail
	detail = GetServerDetail("10.211.55.18")
	logrus.Info(detail)

	detail = GetServerDetail("192.168.116.87")
	logrus.Info(detail)
}
