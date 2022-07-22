package pagination

import (
	"fmt"
	"net/url"
	"strconv"
)

type Pagination struct {
	Limit  int
	Offset int
}

type Error struct {
	Map map[string]interface{}
}

func (e Error) Error() string {
	return fmt.Sprint(e.Map)
}

func (pg Pagination) Validate() error {
	if pg.Limit < 0 {
		return Error{map[string]interface{}{"limit": "Negative values not allowed"}}
	}
	if pg.Offset < 0 {
		return Error{map[string]interface{}{"offset": "Negative values not allowed"}}
	}
	return nil
}

func FromQuery(q url.Values) (Pagination, error) {

	var pg Pagination
	var err error
	errors := make(map[string]interface{}, 2)

	limit := q.Get("limit")
	if limit != "" {
		pg.Limit, err = strconv.Atoi(limit)
		if err != nil {
			errors["limit"] = "Incorrect format. Integer values allowed"
		}
	}
	offset := q.Get("offset")
	if offset != "" {
		pg.Offset, err = strconv.Atoi(offset)
		if err != nil {
			errors["offset"] = "Incorrect format. Integer values allowed"
		}
	}
	if len(errors) > 0 {
		return pg, Error{Map: errors}
	}

	return pg, nil

}
