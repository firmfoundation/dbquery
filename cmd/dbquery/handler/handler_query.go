package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	db "github.com/firmfoundation/dbquery/init"
	"github.com/firmfoundation/dbquery/model"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

var queryTypes map[string]string = map[string]string{"select": "SELECT%", "insert": "INSERT%", "update": "UPDATE%", "delete": "DELETE%"}

func getQueryParams(c *fiber.Ctx) (int, int, string, string, string, error) {
	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil {
		return 0, 0, "", "", "", errors.New("invalid page number")
	}

	pageSize, err := strconv.Atoi(c.Query("page_size", "50"))
	if err != nil {
		return 0, 0, "", "", "", errors.New("invalid page size")
	}

	sort := c.Query("sort", "slowest")
	if sort != "slowest" && sort != "fastest" {
		return 0, 0, "", "", "", errors.New("invalid sort order")
	}

	filterQuery := c.Query("filter_query", "all")
	if filterQuery != "all" && filterQuery != "select" && filterQuery != "insert" && filterQuery != "update" && filterQuery != "delete" {
		return 0, 0, "", "", "", errors.New("invalid filter query")
	}

	dbInstance := c.Query("database_instance_id")
	if dbInstance == "" {
		return 0, 0, "", "", "", errors.New("invalid database instance")
	}

	return page, pageSize, sort, filterQuery, dbInstance, nil
}

func QueryStateHandler(c *fiber.Ctx) error {
	pageInt, pageSizeInt, sorting, filterQuery, dbInstance, err := getQueryParams(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	offset := (pageInt - 1) * pageSizeInt

	var filter string = "*" //default
	if v, ok := queryTypes[strings.ToLower(filterQuery)]; ok {
		filter = v
	}

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
		redis;
		check cache first and query from memory
	*/
	cacheKey := fmt.Sprintf("queries_statistics_%d_%d_%s_%s_%s", pageInt, pageSizeInt, filter, sort, dbInstance)
	cachedResult, err := db.RedisClient.Get(cacheKey).Result()
	if err == nil {
		var raw []map[string]interface{}
		err = json.Unmarshal([]byte(cachedResult), &raw)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusAccepted).JSON(raw)
	}

	/*
		postgreSQL;
		Query data from disk
	*/

	//get db instance
	var instance *gorm.DB
	if v, ok := DbInstanceMap[dbInstance]; ok {
		instance = v
	} else {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "invalid database instance or not connected"})
	}

	queryState := &model.QeryState{}
	result, err := queryState.GetQueryState(instance, pageSizeInt, offset, filter, sort)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	/*
		Cache the results in Redis
		Convert query results to JSON
	*/
	resultJSON, err := json.Marshal(result)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	err = db.RedisClient.Set(cacheKey, resultJSON, 10*time.Second).Err()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusAccepted).JSON(result)
}
