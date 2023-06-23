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
	sort := c.Query("sort")

	pageInt, err := strconv.Atoi(page)
	if err != nil {

	}

	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {

	}

	offset := (pageInt - 1) * pageSizeInt

	queryState := &model.QeryState{}
	result, err := queryState.GetQueryState(db.DB, pageSizeInt, offset)
	if err != nil {
		fmt.Println(err)
	}
	//resultJSON , err = json.Marshal(result)

	return c.Status(200).JSON(result)
}
