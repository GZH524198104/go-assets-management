package model

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"go-asset/store"
	"go.mongodb.org/mongo-driver/bson"
	"math/rand"
	"strconv"
	"time"
)

var rander *rand.Rand

const (
	_REDIS_KEY = "asset:seats"
)

func init() {
	rander = rand.New(rand.NewSource(time.Now().UnixNano()))
}

type Seat struct {
	SeatId string
	X      int
	Y      int
	Weight int
}

func CreateSeatBatch(seats []Seat) error {
	docs := []interface{}{}
	for k := range seats {
		docs = append(docs, bson.D{
			{Key: "_id", Value: seats[k].SeatId},
			{Key: "x", Value: seats[k].X},
			{Key: "y", Value: seats[k].Y},
			{Key: "weight", Value: seats[k].Weight},
		})
	}

	c, err := store.GetMongoClient()
	if err != nil {
		return err
	}
	ctx := context.Background()
	defer c.Disconnect(ctx)

	collection := c.Database("asset").Collection("seats")
	_, err = collection.DeleteMany(ctx, bson.D{})
	if err != nil {
		return err
	}
	_, err = collection.InsertMany(ctx, docs)
	if err != nil {
		return err
	}

	setSeatsCache(_REDIS_KEY, seats, 600)
	return nil
}

func GetAllSeats() ([]Seat, error) {
	seats, _ := GetSeatsCache(_REDIS_KEY)
	if seats == nil {
		c, err := store.GetMongoClient()
		if err != nil {
			return nil, err
		}
		defer c.Disconnect(context.Background())
		collection := c.Database("asset").Collection("seats")

		ctx := context.Background()
		cursor, err := collection.Find(ctx, bson.D{})
		defer cursor.Close(ctx)

		var tmp []Seat
		var result bson.M
		for cursor.Next(ctx) {
			err = cursor.Decode(&result)
			if err != nil {
				return nil, err
			}
			seat := Seat{
				SeatId: result["_id"].(string),
				X:      int(result["x"].(int32)),
				Y:      int(result["y"].(int32)),
				Weight: int(result["weight"].(int32)),
			}
			tmp = append(tmp, seat)
		}
		setSeatsCache(_REDIS_KEY, tmp, 6000)
		return tmp, nil
	}
	return seats, nil
}

func GetSeatsByPersent(persent float64) ([]Seat, string, error) {
	seats, err := GetAllSeats()
	if err != nil {
		return nil, "", err
	}

	var buckets []*Seat
	var total int
	for k := range seats {
		total++
		seat := &seats[k]
		for i := 0; i < seat.Weight; i++ {
			buckets = append(buckets, seat)
		}
	}

	resLength := int(float64(total) * persent / 100)
	var ret []Seat
	var index int
	for i := 0; i < resLength; i++ {
		index = rander.Intn(len(buckets))
		choseOne := buckets[index]
		ret = append(ret, *choseOne)

		j := 0
		for j < len(buckets) {
			if buckets[j] == choseOne {
				buckets = append(buckets[:j], buckets[j+1:]...)
			} else {
				j++
			}
		}
	}

	tmpKey, err := getTmpKey(ret)
	if err != nil {
		return nil, "", err
	}
	return ret, tmpKey, nil
}

func setSeatsCache(key string, seats []Seat, expireTime int) error {
	conn := store.GetRedisConn()
	defer conn.Close()

	buf, _ := json.Marshal(seats)
	_, err := conn.Do("set", key, buf, "ex", expireTime)
	return err
}

func GetSeatsCache(key string) ([]Seat, error) {
	conn := store.GetRedisConn()
	defer conn.Close()

	buf, err := redis.Bytes(conn.Do("get", key))
	if err != nil {
		return nil, err
	}

	var seats []Seat
	err = json.Unmarshal(buf, &seats)
	if err != nil {
		return nil, err
	}
	return seats, nil
}

func getTmpKey(seats []Seat) (string, error) {
	data := []byte(strconv.Itoa(int(time.Now().Unix())))
	has := md5.Sum(data)
	key := fmt.Sprintf("%x", has)

	err := setSeatsCache(key, seats, 6000)
	if err != nil {
		return "", err
	}
	return key, nil
}
