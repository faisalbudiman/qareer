package locations

import (
	"strconv"
	"strings"

	"qareer/pkg/utils"

	sq "github.com/Masterminds/squirrel"
)

var (
	psql      = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	tableName = "locations"
)

type irepo interface {
	Save(Location) error
	Select(filter map[string][]string, loc *[]Location) error
}

type idb interface {
	Save(string, ...interface{}) error
	Select(interface{}, string, ...interface{}) error
}

type repo struct {
	db idb
}

func NewRepo(db idb) repo {
	return repo{
		db: db,
	}
}

func (r repo) Save(loc Location) error {
	sql, args, err := psql.Insert(tableName).Columns("name", "active").Values(loc.Name, loc.Active).ToSql()
	if err != nil {
		return err
	}
	return r.db.Save(sql, args...)
}

func (r repo) Select(filter map[string][]string, loc *[]Location) error {
	query := psql.Select("*").From(tableName).Limit(utils.DefaultLimit).Offset(utils.DefaultSkip)
	if val, ok := filter["name"]; ok {
		for _, v := range val {
			if err := validateFilterName(v); err != nil {
				return err
			}
		}
		if len(val) == 1 {
			query = query.Where("name ILIKE ?", "%"+val[0]+"%")
		} else {
			marks := utils.MapStringToString(val, utils.ReturnQuestionMark)
			values := utils.MapStringToInterface(val, utils.ConvertStringToInterface)
			query = query.Where("name IN ("+strings.Join(marks, ",")+")", values...)
		}
	}

	if val, ok := filter["active"]; ok {
		if len(val) == 1 {
			b, err := strconv.ParseBool(val[0])
			if err != nil {
				return err
			}
			query = query.Where("active = ?", b)
		}
	}

	query, err := utils.SquirrelCommonFilter(query, filter)
	if err != nil {
		return err
	}

	// log.Println("convert query")
	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}

	// log.Println("execute query")
	return r.db.Select(loc, sql, args...)
}
