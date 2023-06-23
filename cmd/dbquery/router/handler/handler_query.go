package handler

import (
	"fmt"
	"strconv"

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

	queryState := &model.QeryState{}
	result, err := queryState.GetQueryState(db.DB, pageSizeInt, offset, sort)
	if err != nil {
		fmt.Println(err)
	}
	//resultJSON , err = json.Marshal(result)

	return c.Status(200).JSON(result)
}
