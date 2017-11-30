package data

import (
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Init mgo and the common DAO

// 数据连接
var Session *mgo.Session
var mgoCloseCh chan bool

// 各个表的Collection对象
var TradeRecords *mgo.Collection
var PlayerUsers *mgo.Collection
var GameRecords *mgo.Collection
var UserRecords *mgo.Collection
var GameRecordsZi *mgo.Collection
var UserRecordsZi *mgo.Collection
var IDGens *mgo.Collection
var WxLogins *mgo.Collection
var LogDiamonds *mgo.Collection
var LogCoins *mgo.Collection
var LogRegists *mgo.Collection
var LogLogins *mgo.Collection
var LogBuildAgencys *mgo.Collection
var LogSetHands *mgo.Collection
var LogOnlines *mgo.Collection
var Agencys *mgo.Collection
var Trends *mgo.Collection
var Notices *mgo.Collection
var Shops *mgo.Collection
var Envs *mgo.Collection
var Boxs *mgo.Collection
var Prizes *mgo.Collection
var Classics *mgo.Collection
var Vips *mgo.Collection
var Mails *mgo.Collection
var DanLists *mgo.Collection
var TaskLists *mgo.Collection
var DanTasks *mgo.Collection
var DanRankings *mgo.Collection
var DanCombats *mgo.Collection
var Seasons *mgo.Collection
var Robots *mgo.Collection
var Bettings *mgo.Collection
var BettingRecords *mgo.Collection
var BettingUsers *mgo.Collection
var Lotterys *mgo.Collection

// 初始化时连接数据库
func InitMgo(dbHost, dbPort, dbUser, dbPassword, dbName string) {
	// get db config from host, port, username, password
	usernameAndPassword := dbUser + ":" + dbPassword + "@"
	if dbUser == "" || dbPassword == "" {
		usernameAndPassword = ""
	}
	if dbPort == "" {
		dbPort = "27017"
	}
	url := "mongodb://" + usernameAndPassword + dbHost + ":" + dbPort + "/" + dbName

	// [mongodb://][user:pass@]host1[:port1][,host2[:port2],...][/database][?options]
	// mongodb://myuser:mypass@localhost:40001,otherhost:40001/mydb
	var err error
	Session, err = mgo.Dial(url)
	if err != nil {
		//return //TODO test
		panic(err)
	}
	go ticker()

	// Optional. Switch the session to a monotonic behavior.
	Session.SetMode(mgo.Monotonic, true)

	// trade_record
	TradeRecords = Session.DB(dbName).C("col_trade_record")
	// user
	PlayerUsers = Session.DB(dbName).C("col_user")
	// record
	GameRecords = Session.DB(dbName).C("col_game_record")
	UserRecords = Session.DB(dbName).C("col_user_record")
	GameRecordsZi = Session.DB(dbName).C("col_game_record_zi")
	UserRecordsZi = Session.DB(dbName).C("col_user_record_zi")
	IDGens = Session.DB(dbName).C("col_id_gen")
	WxLogins = Session.DB(dbName).C("col_wx_login")
	LogDiamonds = Session.DB(dbName).C("col_log_diamond")
	LogCoins = Session.DB(dbName).C("col_log_coin")
	LogRegists = Session.DB(dbName).C("col_log_regist")
	LogLogins = Session.DB(dbName).C("col_log_login")
	LogBuildAgencys = Session.DB(dbName).C("col_log_build_agency")
	LogOnlines = Session.DB(dbName).C("col_log_online")
	LogSetHands = Session.DB(dbName).C("col_log_sethand")
	Agencys = Session.DB(dbName).C("t_user")
	Trends = Session.DB(dbName).C("col_trends")
	Notices = Session.DB(dbName).C("col_notice")
	Shops = Session.DB(dbName).C("col_shop")
	Envs = Session.DB(dbName).C("col_env")
	Boxs = Session.DB(dbName).C("col_box")
	Prizes = Session.DB(dbName).C("col_prize")
	Classics = Session.DB(dbName).C("col_classic")
	Vips = Session.DB(dbName).C("col_vip")
	Mails = Session.DB(dbName).C("col_mail")
	DanLists = Session.DB(dbName).C("col_dan_list")
	TaskLists = Session.DB(dbName).C("col_task_list")
	DanTasks = Session.DB(dbName).C("col_dan_task")
	DanRankings = Session.DB(dbName).C("col_dan_ranking")
	DanCombats = Session.DB(dbName).C("col_dan_combat")
	Seasons = Session.DB(dbName).C("col_dan_season")
	Robots = Session.DB(dbName).C("col_robot")
	Bettings = Session.DB(dbName).C("col_betting")
	BettingRecords = Session.DB(dbName).C("col_betting_record")
	BettingUsers = Session.DB(dbName).C("col_betting_user_record")
	Lotterys = Session.DB(dbName).C("col_lottery")
}

func Close() {
	Session.Close()
	if mgoCloseCh != nil {
		close(mgoCloseCh)
	}
}

// common DAO
// 公用方法

//----------------------

func Insert(collection *mgo.Collection, i interface{}) bool {
	err := collection.Insert(i)
	return Err(err)
}

//----------------------

// 适合一条记录全部更新
func Update(collection *mgo.Collection, query interface{}, i interface{}) bool {
	err := collection.Update(query, i)
	return Err(err)
}
func Upsert(collection *mgo.Collection, query interface{}, i interface{}) bool {
	_, err := collection.Upsert(query, i)
	return Err(err)
}
func UpdateAll(collection *mgo.Collection, query interface{}, i interface{}) bool {
	_, err := collection.UpdateAll(query, i)
	return Err(err)
}
func UpdateByIdAndUserId(collection *mgo.Collection, id, userId string, i interface{}) bool {
	err := collection.Update(GetIdAndUserIdQ(id, userId), i)
	return Err(err)
}

func UpdateByIdAndUserId2(collection *mgo.Collection, id, userId bson.ObjectId, i interface{}) bool {
	err := collection.Update(GetIdAndUserIdBsonQ(id, userId), i)
	return Err(err)
}
func UpdateByIdAndUserIdField(collection *mgo.Collection, id, userId, field string, value interface{}) bool {
	return UpdateByIdAndUserId(collection, id, userId, bson.M{"$set": bson.M{field: value}})
}
func UpdateByIdAndUserIdMap(collection *mgo.Collection, id, userId string, v bson.M) bool {
	return UpdateByIdAndUserId(collection, id, userId, bson.M{"$set": v})
}

func UpdateByIdAndUserIdField2(collection *mgo.Collection, id, userId bson.ObjectId, field string, value interface{}) bool {
	return UpdateByIdAndUserId2(collection, id, userId, bson.M{"$set": bson.M{field: value}})
}
func UpdateByIdAndUserIdMap2(collection *mgo.Collection, id, userId bson.ObjectId, v bson.M) bool {
	return UpdateByIdAndUserId2(collection, id, userId, bson.M{"$set": v})
}

//
func UpdateByQField(collection *mgo.Collection, q interface{}, field string, value interface{}) bool {
	_, err := collection.UpdateAll(q, bson.M{"$set": bson.M{field: value}})
	return Err(err)
}
func UpdateByQI(collection *mgo.Collection, q interface{}, v interface{}) bool {
	_, err := collection.UpdateAll(q, bson.M{"$set": v})
	return Err(err)
}

// 查询条件和值
func UpdateByQMap(collection *mgo.Collection, q interface{}, v interface{}) bool {
	_, err := collection.UpdateAll(q, bson.M{"$set": v})
	return Err(err)
}

//------------------------

// 删除一条
func Delete(collection *mgo.Collection, q interface{}) bool {
	err := collection.Remove(q)
	return Err(err)
}
func DeleteByIdAndUserId(collection *mgo.Collection, id, userId string) bool {
	err := collection.Remove(GetIdAndUserIdQ(id, userId))
	return Err(err)
}
func DeleteByIdAndUserId2(collection *mgo.Collection, id, userId bson.ObjectId) bool {
	err := collection.Remove(GetIdAndUserIdBsonQ(id, userId))
	return Err(err)
}

// 删除所有
func DeleteAllByIdAndUserId(collection *mgo.Collection, id, userId string) bool {
	_, err := collection.RemoveAll(GetIdAndUserIdQ(id, userId))
	return Err(err)
}
func DeleteAllByIdAndUserId2(collection *mgo.Collection, id, userId bson.ObjectId) bool {
	_, err := collection.RemoveAll(GetIdAndUserIdBsonQ(id, userId))
	return Err(err)
}

func DeleteAll(collection *mgo.Collection, q interface{}) bool {
	_, err := collection.RemoveAll(q)
	return Err(err)
}

//-------------------------

func Get(collection *mgo.Collection, id string, i interface{}) {
	collection.FindId(id).One(i)
}
func Get2(collection *mgo.Collection, id bson.ObjectId, i interface{}) {
	collection.FindId(id).One(i)
}
func Get3(collection *mgo.Collection, id string, i interface{}) {
	collection.FindId(bson.ObjectIdHex(id)).One(i)
}

func GetByQ(collection *mgo.Collection, q interface{}, i interface{}) {
	collection.Find(q).One(i)
}
func ListByQ(collection *mgo.Collection, q interface{}, i interface{}) {
	collection.Find(q).All(i)
}

func ListByQLimit(collection *mgo.Collection, q interface{}, i interface{}, limit int) {
	collection.Find(q).Limit(limit).All(i)
}

// 查询某些字段, q是查询条件, fields是字段名列表
func GetByQWithFields(collection *mgo.Collection, q bson.M, fields []string, i interface{}) {
	selector := make(bson.M, len(fields))
	for _, field := range fields {
		selector[field] = true
	}
	collection.Find(q).Select(selector).One(i)
}

// 查询某些字段, q是查询条件, fields是字段名列表
func ListByQWithFields(collection *mgo.Collection, q bson.M, fields []string, i interface{}) {
	selector := make(bson.M, len(fields))
	for _, field := range fields {
		selector[field] = true
	}
	collection.Find(q).Select(selector).All(i)
}
func GetByIdAndUserId(collection *mgo.Collection, id, userId string, i interface{}) {
	collection.Find(GetIdAndUserIdQ(id, userId)).One(i)
}
func GetByIdAndUserId2(collection *mgo.Collection, id, userId bson.ObjectId, i interface{}) {
	collection.Find(GetIdAndUserIdBsonQ(id, userId)).One(i)
}

// 按field去重
func Distinct(collection *mgo.Collection, q bson.M, field string, i interface{}) {
	collection.Find(q).Distinct(field, i)
}

//----------------------

func Count(collection *mgo.Collection, q interface{}) int {
	cnt, err := collection.Find(q).Count()
	if err != nil {
		Err(err)
	}
	return cnt
}

func Has(collection *mgo.Collection, q interface{}) bool {
	if Count(collection, q) > 0 {
		return true
	}
	return false
}

//-----------------

// 得到主键和userId的复合查询条件
func GetIdAndUserIdQ(id, userId string) bson.M {
	return bson.M{"_id": bson.ObjectIdHex(id), "UserId": bson.ObjectIdHex(userId)}
}
func GetIdAndUserIdBsonQ(id, userId bson.ObjectId) bson.M {
	return bson.M{"_id": id, "UserId": userId}
}

// DB处理错误
func Err(err error) bool {
	if err != nil {
		//fmt.Println(err)
		// 删除时, 查找
		if err.Error() == "not found" {
			return true
		}
		return false
	}
	return true
}

// 检查mognodb是否lost connection
// 每个请求之前都要检查!!
func CheckMongoSessionLost() {
	// fmt.Println("检查CheckMongoSessionLostErr")
	err := Session.Ping()
	if err != nil {
		//Log("Lost connection to db!")
		Session.Refresh()
		err = Session.Ping()
		if err == nil {
			//Log("Reconnect to db successful.")
		} else {
			//Log("重连失败!!!! 警告")
		}
	}
}

//计时器
func ticker() {
	mgoCloseCh = make(chan bool, 1)
	tick := time.Tick(time.Minute)
	for {
		select {
		case <-tick:
			CheckMongoSessionLost()
		case <-mgoCloseCh:
			return
		}
	}
}

// 分页, 排序处理
func parsePageAndSort(pageNumber, pageSize int, sortField string, isAsc bool) (skipNum int, sortFieldR string) {
	skipNum = (pageNumber - 1) * pageSize
	if sortField == "" {
		sortField = "UpdatedTime"
	}
	if !isAsc {
		sortFieldR = "-" + sortField
	} else {
		sortFieldR = sortField
	}
	return
}
