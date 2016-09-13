package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/garyburd/redigo/redis"
)

func getCurrentPath() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}
func callredis_updateEXE(filePath string, filename string, argA, argB string) {
	//往redis里插入数据之前先清空一下redis
	client := connRedis()
	defer client.Close()

	_, err := client.Do("flushall")
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Current path:", getCurrentPath())
	cmd := exec.Command(filePath+filename, argA, argB)
	cmd.Dir = filePath
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		fmt.Println("Run error:", err.Error())
		panic(err.Error())
	}
	fmt.Println("Output info:", out.String())
}
func dataConverStrings(res interface{}, err error) []string {
	if err != nil {
		panic(err.Error())
	}
	value, err := redis.Strings(res, err)
	if err != nil {
		panic(err.Error())
	}
	return value
}

func dataConverString(res interface{}, err error) string {
	if err != nil {
		panic(err.Error())
	}
	value, err := redis.String(res, err)
	if err != nil {
		panic(err.Error())
	}
	return value
}
func redisDoStrings(client redis.Conn, args ...string) []string {
	switch len(args) {
	case 1:
		res, err := client.Do(args[0])
		value := dataConverStrings(res, err)
		return value
	case 2:
		res, err := client.Do(args[0], args[1])
		value := dataConverStrings(res, err)
		return value
	case 3:
		res, err := client.Do(args[0], args[1], args[2])
		value := dataConverStrings(res, err)
		return value
	default:
		fmt.Println("not args into")
		panic("unrecognized escape character")
	}

}

func redisDoString(client redis.Conn, args ...string) string {
	switch len(args) {
	case 1:
		res, err := client.Do(args[0])
		value := dataConverString(res, err)
		return value
	case 2:
		res, err := client.Do(args[0], args[1])
		value := dataConverString(res, err)
		return value
	case 3:
		res, err := client.Do(args[0], args[1], args[2])
		value := dataConverString(res, err)
		return value
	default:
		fmt.Println("not args into")
		panic("unrecognized escape character")
	}

}
func operateRedis(fileName string) {
	file, ERR := os.Create(fileName)
	if ERR != nil {
		panic(ERR.Error())
	}
	defer file.Close()

	client := connRedis()
	defer client.Close()

	keys := redisDoStrings(client, "KEYS", "*")
	for i, k := range keys {
		fmt.Println("keys:", i, k)
		_, err := client.Do("SADD", "redisOld", k)
		if err != nil {
			panic(err.Error())
		}

	}
	redisOld := redisDoStrings(client, "SORT", "redisOld", "ALPHA")
	for _, k := range redisOld {
		//fmt.Println("new Keys:", i, k)
		keytype := redisDoString(client, "type", k)

		if keytype == "string" {
			value := redisDoString(client, "GET", k)
			FileWrite(file, "\n")
			FileWrite(file, k+":"+value)
		} else if keytype == "hash" {
			value := redisDoStrings(client, "HKEYS", k)
			hashkey := SortHashKey(value, client)
			fmt.Println(hashkey)
			WriteHashValuetoFile(k, hashkey, file, client)
		}
	}

}
func FileWrite(file *os.File, str string) {
	_, err := file.WriteString(str)
	if err != nil {
		panic(err.Error())
	}
}
func SortHashKey(key []string, client redis.Conn) []string {
	for _, k := range key {
		_, err := client.Do("SADD", "HashKey", k)
		if err != nil {
			panic(err.Error())
		}
	}
	value := redisDoStrings(client, "SORT", "HashKey", "ALPHA")
	return value
}
func WriteHashValuetoFile(hash string, keys []string, file *os.File, client redis.Conn) {
	FileWrite(file, "\n")
	FileWrite(file, hash+":\n")
	for _, k := range keys {
		value := redisDoString(client, "HGET", hash, k)
		FileWrite(file, k+"    ")
		FileWrite(file, value+"\n")
	}
	_, err := client.Do("DEL", "HashKey")
	if err != nil {
		panic(err.Error())
	}
}
func connRedis() redis.Conn {
	client, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		fmt.Println("Connect redis error:", err.Error())
		panic(err.Error())
	}

	return client
}
func main() {

	filePath := getCurrentPath() + "/old/"
	callredis_updateEXE(filePath, "redis_update.exe", "-log", "true")
	operateRedis("OldData.txt")

	filePath = getCurrentPath() + "/new/"
	callredis_updateEXE(filePath, "/redis_update.exe", "-log", "true")
	operateRedis("NewData.txt")
}
