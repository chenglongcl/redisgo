package redisgo

import (
	"reflect"
	"testing"
	"time"
)

type User struct {
	Name string
	Age  int
}

func NoError(t *testing.T, err error) {
	if err != nil {
		t.Error(err)
	}
}

func Error(t *testing.T, err error) {
	if err == nil {
		t.Error("Expected an error.")
	}
}

func Equal(t *testing.T, expected, actual interface{}) {
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Not equal: \n"+
			"expected: %+v\n"+
			"actual  : %+v", expected, actual)
	}
}

func getCacher() *Cacher {
	c, err := New(
		Options{
			Network:  "tcp",
			Addr:     "127.0.0.1:6379",
			Password: "20190604",
			Prefix:   "zengate_",
		})
	if err != nil {
		panic(err)
	}
	return c
}

func TestGetSet(t *testing.T) {
	var err error
	c := getCacher()

	// int
	err = c.Set("age", "23", 30)
	NoError(t, err)
	valInt, err := c.GetInt("age")
	NoError(t, err)
	Equal(t, 23, valInt)

	// string
	err = c.Set("name", "corel", 30)
	NoError(t, err)
	valString, err := c.GetString("name")
	NoError(t, err)
	Equal(t, "corel", valString)

	// bool
	err = c.Set("subscribe", true, 30)
	NoError(t, err)
	valBool, err := c.GetBool("subscribe")
	NoError(t, err)
	Equal(t, true, valBool)

	// user
	user := &User{
		Name: "corel",
		Age:  23,
	}
	err = c.Set("user", user, 30)
	NoError(t, err)
	valUser := &User{}
	err = c.GetObject("user", valUser)
	NoError(t, err)
	Equal(t, "corel", valUser.Name)
	Equal(t, 23, valUser.Age)
}

func TestIncrDecr(t *testing.T) {
	var err error
	c := getCacher()
	c.Del("seq")
	val, err := c.Incr("seq")
	NoError(t, err)
	Equal(t, int64(1), val)
	val, err = c.Incr("seq")
	NoError(t, err)
	Equal(t, int64(2), val)
	val, err = c.IncrBy("seq", 5)
	NoError(t, err)
	Equal(t, int64(7), val)
	val, err = c.Decr("seq")
	NoError(t, err)
	Equal(t, int64(6), val)
	val, err = c.DecrBy("seq", 5)
	NoError(t, err)
	Equal(t, int64(1), val)
}

func TestHKeys(t *testing.T) {
	var err error
	c := getCacher()
	c.Del("hKeyTest")
	_, err = c.HSet("hKeyTest", "field1", "foo")
	NoError(t, err)
	_, err = c.HSet("hKeyTest", "field2", "bar")
	NoError(t, err)
	val, err := c.HKeys("hKeyTest")
	NoError(t, err)
	Equal(t, val, []string{"field1", "field2"})
}

func TestExpire(t *testing.T) {
	var err error
	c := getCacher()
	err = c.Set("name", "corel", 1)
	NoError(t, err)

	time.Sleep(2 * time.Second)

	_, err = c.GetString("name")
	Error(t, err)
}

func TestHash(t *testing.T) {
	var err error
	c := getCacher()
	m := make(map[string]interface{})
	m["name"] = "corel"
	m["age"] = 23
	err = c.HMSet("huser", m, 10)
	NoError(t, err)

	age, err := c.HGetInt("huser", "age")
	NoError(t, err)
	Equal(t, m["age"], age)
}

func TestSortedSet(t *testing.T) {
	var err error
	c := getCacher()
	_, err = c.ZAdd("scores", 82, "corel")
	NoError(t, err)
	_, err = c.ZAdd("scores", 86, "zen")
	NoError(t, err)
	score, err := c.ZScore("scores", "corel")
	NoError(t, err)
	Equal(t, int64(82), score)
}

func TestHDel(t *testing.T) {
	var err error
	c := getCacher()
	_, err = c.HSet("huser", "1", "haha")
	NoError(t, err)
	_, err = c.HDel("huser", "1")
	NoError(t, err)
}

func TestHIncrby(t *testing.T) {
	var err error
	c := getCacher()
	_, err = c.HIncrby("hnum", "num", 1)
	NoError(t, err)
}

func TestZIncrby(t *testing.T) {
	var err error
	c := getCacher()
	_, err = c.ZIncrby("hNum", "num", 1)
	NoError(t, err)
}

func TestZcard(t *testing.T) {
	var err error
	c := getCacher()
	val, err := c.ZCard("hNum")
	Equal(t, int(1), val)
	NoError(t, err)
}

func TestType(t *testing.T) {
	c := getCacher()
	var err error
	keyString := "typeTestString"
	keyList := "typeTestList"
	keyHash := "typeTestHash"
	keyZset := "typeTestZset"
	keyNone := "typeTestNone"

	// Clean up keys before starting
	c.Del(keyString)
	c.Del(keyList)
	c.Del(keyHash)
	c.Del(keyZset)
	c.Del(keyNone)

	// Test String
	err = c.Set(keyString, "hello", 30)
	NoError(t, err)
	typeVal, err := c.Type(keyString)
	NoError(t, err)
	Equal(t, "string", typeVal)

	// Test List
	err = c.LPush(keyList, "world")
	NoError(t, err)
	typeVal, err = c.Type(keyList)
	NoError(t, err)
	Equal(t, "list", typeVal)

	// Test Hash
	_, err = c.HSet(keyHash, "field", "value")
	NoError(t, err)
	typeVal, err = c.Type(keyHash)
	NoError(t, err)
	Equal(t, "hash", typeVal)

	// Test ZSet (Sorted Set)
	_, err = c.ZAdd(keyZset, 1, "member1")
	NoError(t, err)
	typeVal, err = c.Type(keyZset)
	NoError(t, err)
	Equal(t, "zset", typeVal)

	// Test None (non-existent key)
	typeVal, err = c.Type(keyNone)
	NoError(t, err)
	Equal(t, "none", typeVal)

	// Clean up after test
	c.Del(keyString)
	c.Del(keyList)
	c.Del(keyHash)
	c.Del(keyZset)
}
