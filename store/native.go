package store

import (
	"context"
	"database/sql"
	"dna-test/config"
	"dna-test/models"
	"errors"
	"fmt"
	"github.com/georgysavva/scany/v2/dbscan"
	_ "github.com/lib/pq"
	"reflect"
	"strings"
)

type DbClientNative struct {
	Db *sql.DB
}

func NewNativeConn(db *sql.DB) *DbClientNative {
	return &DbClientNative{db}
}

func NewNativeDB(config *config.Config, out *models.DbClient) error {
	connectionString := fmt.Sprintf("host=%v port=%v user=%v dbname=%v password=%v sslmode=%v",
		config.DbHost, config.DbPort, config.DbUser, config.DbName, config.DbPass, config.DbMode)

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		panic(err)
	}

	conn := NewNativeConn(db)
	if out != nil && db != nil {
		*out = conn
	}
	return nil
}

func (c *DbClientNative) GetTable(ctx context.Context, table string, query models.FilterParams, out interface{}) error {
	if out == nil {
		return errors.New("cannot assign values into nil")
	}

	isSlice := false
	header := reflect.TypeOf(out)
	if header.Kind() != reflect.Ptr {
		return errors.New("should be pointer to the struct, not the struct")
	}
	header = header.Elem()
	if header.Kind() == reflect.Slice {
		header = header.Elem()
		isSlice = true
	}
	var row []string
	ConstructRow(header, &row)
	fields := strings.Join(row, ",")
	sqlQuery := fmt.Sprintf(`select %v from %v`, fields, table)
	if len(query.Filter) != 0 {
		sqlQuery = fmt.Sprintf(`%v where %v`, sqlQuery, query.Filter)
	}
	if len(query.Orderings) != 0 {
		sqlQuery = fmt.Sprintf(`%v order by %v`, sqlQuery, query.Orderings)
	}

	rows, err := c.Db.QueryContext(ctx, sqlQuery)
	if err != nil {
		return err
	}
	defer rows.Close()

	return ScanRows(rows, out, isSlice)
}

func (c *DbClientNative) Select(ctx context.Context, table string, params models.FilterParams, out interface{}) error {
	if out == nil {
		return errors.New("cannot assign values into nil")
	}

	isSlice := false
	header := reflect.TypeOf(out)
	if header.Kind() != reflect.Ptr {
		return errors.New("should be pointer to the struct, not the struct")
	}
	header = header.Elem()
	if header.Kind() == reflect.Slice {
		isSlice = true
	}
	sqlQuery := fmt.Sprintf(`select %v from %v`, params.Select, table)
	if len(params.Filter) != 0 {
		sqlQuery = fmt.Sprintf(`%v where %v`, sqlQuery, params.Filter)
	}
	if len(params.Orderings) != 0 {
		sqlQuery = fmt.Sprintf(`%v order by %v`, sqlQuery, params.Orderings)
	}

	rows, err := c.Db.QueryContext(ctx, sqlQuery)
	if err != nil {
		return err
	}
	defer rows.Close()

	return ScanRows(rows, out, isSlice)
}

func (c *DbClientNative) PingClient(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (c *DbClientNative) Create(ctx context.Context, input interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (c *DbClientNative) Get(ctx context.Context, params models.FilterParams, out interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (c *DbClientNative) Update(ctx context.Context, params models.FilterParams, input interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (c *DbClientNative) Upsert(ctx context.Context, params models.FilterParams, input interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (c *DbClientNative) Delete(ctx context.Context, params models.FilterParams, input interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (c *DbClientNative) CloseClient() error {
	//TODO implement me
	panic("implement me")
}

func ConstructRow(header reflect.Type, row *[]string) {
	var field reflect.StructField
	numFields := header.NumField()
	for i := 0; i < numFields; i++ {
		field = header.Field(i)
		if alias, ok := field.Tag.Lookup("db"); ok && alias != "-" {
			if field.Type.Kind() == reflect.Struct && alias == "" {
				ConstructRow(field.Type, row)
			} else {
				*row = append(*row, alias)
			}
		}
	}
}

func ScanRows(rows *sql.Rows, out interface{}, isSlice bool) error {
	if isSlice {
		err := dbscan.ScanAll(out, rows)
		if err != nil {
			return err
		}
		s := reflect.ValueOf(out)
		if s.Elem().IsNil() {
			return errors.New(models.DB_ERROR_NOT_FOUND)
		}
	}
	return dbscan.ScanOne(out, rows)
}
