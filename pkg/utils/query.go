package utils

import (
	"strconv"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"
)

const (
	limitName           = "limit"
	skipName            = "skip"
	DefaultLimit uint64 = 10
	maxLimit     uint64 = 30
	DefaultSkip  uint64 = 0
)

func SquirrelFilterByID(query sq.SelectBuilder, filter map[string][]string) (sq.SelectBuilder, error) {
	if val, ok := filter["id"]; ok {
		for _, v := range val {
			_, err := strconv.ParseUint(v, 10, 32)
			if err != nil {
				return query, err
			}
		}
		if len(val) == 1 {
			query = query.Where("id = ?", val[0])
		} else {
			marks := MapStringToString(val, ReturnQuestionMark)
			values := MapStringToInterface(val, ConvertStringToInterface)
			query = query.Where("id IN ("+strings.Join(marks, ",")+")", values...)
		}
	}

	return query, nil
}

func SquirrelFilterByCreatedAT(query sq.SelectBuilder, filter map[string][]string) (sq.SelectBuilder, error) {
	if val, ok := filter["created_at"]; ok {
		for _, v := range val {
			_, err := time.Parse("2006-01-02 03:04:05", v)
			if err != nil {
				return query, err
			}
		}
		if len(val) == 2 {
			query = query.Where("created_at between ? and ?", val[0], val[1])
		}
	}

	return query, nil
}

func SquirrelFilterByUpdatedAT(query sq.SelectBuilder, filter map[string][]string) (sq.SelectBuilder, error) {
	if val, ok := filter["updated_at"]; ok {
		for _, v := range val {
			_, err := time.Parse("2006-01-02 03:04:05", v)
			if err != nil {
				return query, err
			}
		}
		if len(val) == 2 {
			query = query.Where("updated_at between ? and ?", val[0], val[1])
		}
	}

	return query, nil
}

func SquirrelLimit(query sq.SelectBuilder, filter map[string][]string) (sq.SelectBuilder, error) {
	if val, ok := filter[limitName]; ok {
		if len(val) == 1 {
			i, err := strconv.ParseUint(val[0], 10, 32)
			if err != nil {
				return query, err
			}
			if i <= maxLimit {
				query = query.Limit(i)
			} else {
				query = query.Limit(maxLimit)
			}

		}
	}
	return query, nil
}

func SquirrelSkip(query sq.SelectBuilder, filter map[string][]string) (sq.SelectBuilder, error) {
	// log.Println("check skip")
	if val, ok := filter[skipName]; ok {
		if len(val) == 1 {
			i, err := strconv.ParseUint(val[0], 10, 32)
			if err != nil {
				return query, err
			}
			query = query.Offset(i)
		}
	}
	return query, nil
}

func SquirrelCommonFilter(query sq.SelectBuilder, filter map[string][]string) (sq.SelectBuilder, error) {
	query, err := SquirrelFilterByID(query, filter)
	if err != nil {
		return query, err
	}

	query, err = SquirrelFilterByCreatedAT(query, filter)
	if err != nil {
		return query, err
	}

	query, err = SquirrelFilterByUpdatedAT(query, filter)
	if err != nil {
		return query, err
	}

	query, err = SquirrelLimit(query, filter)
	if err != nil {
		return query, err
	}

	query, err = SquirrelSkip(query, filter)
	if err != nil {
		return query, err
	}

	return query, nil
}
