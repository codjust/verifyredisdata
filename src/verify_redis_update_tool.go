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
func callredis_updateEXE(filePath string, argA, argB string) {

	//往redis里插入数据之前先清空一下redis
	client := connRedis()
	defer client.Close()

	_, err := client.Do("flushall")
	if err != nil {
		panic(err.Error())
	}

	//	fmt.Println("Current path:", getCurrentPath())
	cmd := exec.Command(filePath, argA, argB)

	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		fmt.Println("Run error:", err.Error())
		panic(err.Error())
	}

	fmt.Println("Output info:", out.String())

}

func operateRedis() {
	client := connRedis()
	defer client.Close()

	res, err := client.Do("KEYS", "*")
	if err != nil {
		fmt.Println("operate redis error:", err.Error())
		panic(err.Error())
	}

	keys, err := redis.Strings(res, err)
	if err != nil {
		panic(err.Error())
	}

	for i, k := range keys {
		fmt.Println("keys:", i, k)
		_, err := client.Do("SADD", "redisOld", k)
		if err != nil {
			panic(err.Error())
		}
	}

	res, err = client.Do("SORT", "redisOld", "ALPHA")
	if err != nil {
		panic(err.Error())
	}
	redisOld, err := redis.Strings(res, err)
	if err != nil {
		panic(err.Error())
	}
	for i, k := range redisOld {
		fmt.Println("new Keys:", i, k)
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

	callredis_updateEXE("./new/redis_update.exe", "-log", "true")
	operateRedis()
}
