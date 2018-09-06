package connection

import (
	"fmt"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/globalsign/mgo"
	"github.com/kr/pretty"
)

func TestInfo_clone(t *testing.T) {
	info := &Info{
		Username: "user",
		Password: "pwd",
		Host:     "host1",
		Port:     "1",
	}

	info2 := info.clone("", "")
	if !reflect.DeepEqual(info, info2) {
		pretty.Ldiff(t, info, info2)
		t.Error("Bad info clone2")
	}

	info3 := info.clone("host3", "")
	info.Host = "host3"
	if !reflect.DeepEqual(info, info3) {
		pretty.Ldiff(t, info, info3)
		t.Error("Bad info clone2")
	}

	info4 := info.clone("host4", "4")
	info.Host = "host4"
	info.Port = "4"
	if !reflect.DeepEqual(info, info4) {
		pretty.Ldiff(t, info, info4)
		t.Error("Bad info clone2")
	}

	info5 := info.clone("", "5")
	info.Port = "5"
	if !reflect.DeepEqual(info, info5) {
		pretty.Ldiff(t, info, info5)
		t.Error("Bad info clone2")
	}
}

func TestInfo_CreateSession(t *testing.T) {
	info := &Info{
		Username:              "",
		Password:              "",
		Host:                  "localhost",
		Port:                  "27017",
		AuthSource:            "admin",
		Ssl:                   true,
		SslCaCerts:            "test",
		SslInsecureSkipVerify: false,
	}

	_, err := info.CreateSession()
	if err == nil {
		t.Error("Expected connection to fail")
	}

}

func TestInfo_generateDialInfo(t *testing.T) {
	info := &Info{
		Host:       "localhost",
		Port:       "27017",
		AuthSource: "admin",
	}
	dialInfo := info.generateDialInfo()

	expectedDialInfo := &mgo.DialInfo{
		Addrs:          []string{"localhost:27017"},
		Username:       "",
		Password:       "",
		Source:         "admin",
		Direct:         true,
		FailFast:       true,
		Timeout:        time.Duration(10) * time.Second,
		PoolTimeout:    time.Duration(10) * time.Second,
		ReadTimeout:    time.Duration(10) * time.Second,
		ReadPreference: &mgo.ReadPreference{Mode: mgo.PrimaryPreferred},
	}

	if !reflect.DeepEqual(dialInfo, expectedDialInfo) {
		fmt.Println(pretty.Diff(dialInfo, expectedDialInfo))
		t.Error("Bad dial info")
	}
}

func Test_addSSL(t *testing.T) {
	dialInfo := &mgo.DialInfo{
		Addrs:       []string{"localhost"},
		Username:    "",
		Password:    "",
		Source:      "admin",
		FailFast:    true,
		Timeout:     time.Duration(1) * time.Second,
		PoolTimeout: time.Duration(1) * time.Second,
		ReadTimeout: time.Duration(1) * time.Second,
	}

	addSSL(dialInfo, false, "")

	if dialInfo.DialServer == nil {
		t.Error("Nil dialServer")
	}
}

func Test_addSSL_EmptyPEM(t *testing.T) {
	dialInfo := &mgo.DialInfo{
		Addrs:       []string{"localhost"},
		Username:    "",
		Password:    "",
		Source:      "admin",
		FailFast:    true,
		Timeout:     time.Duration(1) * time.Second,
		PoolTimeout: time.Duration(1) * time.Second,
		ReadTimeout: time.Duration(1) * time.Second,
	}

	addSSL(dialInfo, false, filepath.Join("testdata", "empty.pem"))

	if dialInfo.DialServer == nil {
		t.Error("Nil dialServer")
	}
}
