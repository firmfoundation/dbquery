package handler

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	db "github.com/firmfoundation/dbquery/init"
	"github.com/firmfoundation/dbquery/model"
	"github.com/gofiber/fiber/v2"
)

func QueryStateHandler(c *fiber.Ctx) error {
	page := c.Query("page")
	pageSize := c.Query("page_size")
	sorting := c.Query("sort")

	pageInt, err := strconv.Atoi(page)
	if err != nil {

	}

	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {

	}

	offset := (pageInt - 1) * pageSizeInt

	/*
		default sort is DESC
		slowest = quries taking longer execution time fetched first (slowest queries), DESC
		fastest = quriest executied fast feched first, ASC
	*/
	var sort string = "DESC"
	if sorting == "slowest" {
		sort = "DESC"
	} else if sorting == "fastest" {
		sort = "ASC"
	}

	/*
		check cache first
	*/

	cacheKey := fmt.Sprintf("queries_statistics_%d_%d_%s", pageInt, pageSizeInt, sort)
	cachedResult, err := db.RedisClient.Get(cacheKey).Result()
	if err == nil {
		var raw []map[string]interface{}
		err = json.Unmarshal([]byte(cachedResult), &raw)
		if err != nil {
			panic(err)
		}
		return c.Status(200).JSON(raw)
	}

	/*
		Query data from postgreSQL
	*/
	queryState := &model.QeryState{}
	result, err := queryState.GetQueryState(db.DB, pageSizeInt, offset, sort)
	if err != nil {
		fmt.Println(err)
	}

	/*
		Cache the results in Redis
		Convert query results to JSON
	*/
	resultJSON, err := json.Marshal(result)
	if err != nil {
		// Handle JSON marshaling error
	}

	err = db.RedisClient.Set(cacheKey, resultJSON, 10*time.Second).Err()
	if err != nil {
		// Handle cache set error
	}

	return c.Status(200).JSON(result)
}