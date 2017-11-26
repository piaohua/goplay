//微信登录
package data

import (
	"testing"
	"utils"

	"gopkg.in/mgo.v2/bson"
	//"gopkg.in/mgo.v2"
)

/*
func TestWxLogin(t *testing.T) {
	collection := "col_wx_login"
	_DBNAME := "test"
	host := "127.0.0.1:2225"
	var se *mgo.Session
	var err error
	se, err = mgo.Dial(host)
	se.DB(_DBNAME).C(collection)

	openid := "o48K-wCLzOgvD7b2_kbllFcNHDmQ"
	d := &WXLogin{
		OpenId: openid,
		//AccessToken: openid,
	}
	//err = d.Get()
	//err = d.GetByToken()
	//err = d.Save()
	err = d.Update()
	t.Log(d, err)

	//u := &User{
	//	Userid: "16007",
	//}
	//err = u.Get()
	//t.Log(u, err)
}
*/

func TestRecord(t *testing.T) {
	InitMgo()
	list, err := GetRecords("10026", 1)
	t.Log(err)
	t.Logf("%#v", list)
	for k, v := range list {
		t.Log("k ", k, v.Roomid, v.Ctime)
	}
}

func TestGroup(t *testing.T) {
	InitMgo()
	var types []string
	ListByQ(UserRecords, bson.M{"$group": bson.M{"userid": "$userid"}}, &types)
	t.Logf("%#v", types)
	m := bson.M{
		"$match": bson.M{
		//"diamond": bson.M{"$ne": 0},
		},
	}
	o := bson.M{
		"$project": bson.M{
			"_id":     1,
			"diamond": 1,
		},
	}
	n := bson.M{
		"$group": bson.M{
			"_id": "$_id",
		},
	}
	operations := []bson.M{m, o, n}
	result := []bson.M{}
	pipe := PlayerUsers.Pipe(operations)
	err2 := pipe.All(&result)
	t.Logf("%#v", err2)
	t.Logf("%#v", result)
}

func TestOr(t *testing.T) {
	InitMgo()
	or := []bson.M{bson.M{"userid": "10030"}}
	m := bson.M{
		"$match": bson.M{
			"$or": or,
		},
	}
	operations := []bson.M{m}
	result := []bson.M{}
	pipe := TradeRecords.Pipe(operations)
	err2 := pipe.All(&result)
	t.Logf("%#v", err2)
	t.Logf("%#v", result)
}

func TestID(t *testing.T) {
	//dayStamp := utils.Stamp2Time(utils.TimestampYesterday())
	dayStamp := utils.Stamp2Time(utils.TimestampToday())
	id := bson.NewObjectIdWithTime(dayStamp).Hex()
	t.Logf("%s", id)
	id2 := bson.NewObjectIdWithTime(dayStamp).String()
	t.Logf("%s", id2)
	t.Logf("%#v", dayStamp)
}
