package main

import (
	"fmt"
	"log"

	"github.com/gomodule/redigo/redis"
)

type Mahasiswa struct {
	Name      string  `redis:"name"`
	StudentId string  `redis:"student_id"`
	GPA       float64 `redis:"gpa"`
	Semester  int     `redis:"semester"`
}

func main() {
	// Create Connection (Not Secure for Concurrency)
	conn, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		log.Panic(err)
	}

	// Send Redis Command
	_, err = conn.Do("HSET", "student:1", "name", "Andromeda", "student_id", "12345", "gpa", "3.45", "semester", "4")
	if err != nil {
		log.Panic(err)
	}

	// Type Convertion
	name, err := redis.String(conn.Do("HGET", "student:1", "name"))
	if err != nil {
		log.Panic(err)
	}

	studentId, err := redis.String(conn.Do("HGET", "student:1", "student_id"))
	if err != nil {
		log.Panic(err)
	}

	gpa, err := redis.Float64(conn.Do("HGET", "student:1", "gpa"))
	if err != nil {
		log.Panic(err)
	}

	semester, err := redis.Int(conn.Do("HGET", "student:1", "semester"))
	if err != nil {
		log.Panic(err)
	}

	fmt.Println(name)      // Andromeda
	fmt.Println(studentId) // 12345
	fmt.Println(gpa)       // 3.45
	fmt.Println(semester)  // 4

	// Get All Data Hash/Map
	replyMap, err := redis.StringMap(conn.Do("HGETALL", "student:1"))
	if err != nil {
		log.Panic(err)
	}
	fmt.Println(replyMap) // map[gpa:3.45 name:Andromeda name:Andromeda student_id:12345 semester:4]

	replyAny, err := redis.Values(conn.Do("HGETALL", "student:1"))
	if err != nil {
		log.Panic(err)
	}

	// Save All Data Hash/Map to Struct
	var mahasiswa Mahasiswa

	err = redis.ScanStruct(replyAny, &mahasiswa)
	if err != nil {
		log.Panic(err)
	}

	fmt.Printf("%+v\n", mahasiswa) // {Name:Andromeda StudentId:12345 GPA:3.45 Semester:4}

	// Create Connection Polling (Secure for Concurrency)
	pool := &redis.Pool{
		MaxIdle: 10,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "localhost:6379")
		},
		MaxActive: 10,
	}

	conn = pool.Get()
	defer conn.Close()

	reply, err := conn.Do("SET", "test", "oke")
	if err != nil {
		log.Panic(err)
	}

	fmt.Println(reply) // OK
}
