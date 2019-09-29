package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-asset/model"
	"io/ioutil"
	"strconv"
	"strings"
)

const (
	host = "localhost:8080"
)

func PostSeatsCsv(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(500, gin.H{
			"errMsg": err.Error(),
		})
		return
	}

	f, err := file.Open()
	if err != nil {
		c.JSON(500, gin.H{
			"errMsg": err.Error(),
		})
		return
	}

	all, err := ioutil.ReadAll(f)
	if err != nil {
		c.JSON(500, gin.H{
			"errMsg": err.Error(),
		})
		return
	}

	arrStr := strings.Split(string(all), "\n")
	var seats []model.Seat
	for k := range arrStr {
		line := arrStr[k]
		if len(line) <= 0 {
			continue
		}
		arrItem := strings.Split(strings.TrimSpace(line), ",")
		seatId := arrItem[0]
		x, _ := strconv.Atoi(arrItem[1])
		y, _ := strconv.Atoi(arrItem[2])
		weight, _ := strconv.Atoi(arrItem[3])
		seats = append(seats, model.Seat{
			SeatId: seatId,
			X:      x,
			Y:      y,
			Weight: weight,
		})
	}

	err = model.CreateSeatBatch(seats)
	if err != nil {
		c.JSON(500, gin.H{
			"errMsg": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"msg": "success",
	})
}

func GetSeatsByPersent(c *gin.Context) {
	persent := c.Param("persent")
	p, err := strconv.ParseFloat(persent, 64)
	if err != nil || p < 0 || p > 100 {
		var errMsg string
		if err == nil {
			errMsg = fmt.Sprintf("Illegal persion value:%f", p)
		} else {
			errMsg = err.Error()
		}
		c.JSON(400, gin.H{
			"errMsg": errMsg,
		})
		return
	}

	seats, key, err := model.GetSeatsByPersent(p)
	if err != nil {
		c.JSON(400, gin.H{
			"errMsg": err.Error(),
		})
		return
	}

	var ids []string

	for k := range seats {
		ids = append(ids, seats[k].SeatId)
	}

	c.JSON(200, gin.H{
		"ids":      ids,
		"iamgeUrl": "https://" + host + "/image/" + key,
	})
}

func GetBestPathImage(c *gin.Context) {
	key := c.Param("key")
	seats, err := model.GetSeatsCache(key)
	if err != nil {
		c.JSON(500, gin.H{
			"errMsg": err.Error(),
		})
		return
	}
	x := 0
	y := 0
	var min int
	var next *model.Seat
	var nextIndex int
	var path []model.Seat

	for len(seats) > 0 {
		min = 9999999
		for k := range seats {
			distanct := getDistance(x, y, seats[k].X, seats[k].Y)
			if distanct < min {
				next = &seats[k]
				nextIndex = k
				min = distanct
			}
		}
		path = append(path, *next)
		seats = append(seats[:nextIndex], seats[nextIndex+1:]...)
	}

	allSeats, err := model.GetAllSeats()
	printer := &model.SeatsPrinter{
		Seats:  allSeats,
		Path:   path,
		Weight: 600,
		Height: 600,
		Space:  50,
	}

	data, err := printer.PrintBestPath()
	if err != nil {
		c.JSON(500, gin.H{
			"errMsg": err.Error(),
		})
		return
	}

	c.Header("Content-Disposition", "attachment; filename=\"path.png\"")
	c.Data(200, "application/octet-stream", data)
}

func getDistance(x int, y int, x1 int, y1 int) int {
	return (x-x1)*(x-x1) + (y-y1)*(y-y1)
}
