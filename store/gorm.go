package store

import (
	"context"
	"dna-test/config"
	"dna-test/models"
	"errors"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DbClientGorm struct {
	Db *gorm.DB
}

func NewGormConn(db *gorm.DB) *DbClientGorm {
	return &DbClientGorm{db}
}

func NewGormDB(config *config.Config, out *models.DbClient) error {
	logLevel := logger.Silent
	if config.DbLogMode {
		logLevel = logger.Info
	}

	connectionString := fmt.Sprintf("host=%v port=%v user=%v dbname=%v password=%v sslmode=%v",
		config.DbHost, config.DbPort, config.DbUser, config.DbName, config.DbPass, config.DbMode)

	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return err
	}

	conn := NewGormConn(db)
	if out != nil && db != nil {
		*out = conn
	}
	return nil
}

func (c *DbClientGorm) CloseClient() error {
	sqlDb, err := c.Db.DB()
	if err != nil {
		return err
	}
	return sqlDb.Close()
}

func (c *DbClientGorm) PingClient(ctx context.Context) error {
	sqlDB, err := c.Db.DB()
	if err != nil {
		return err
	}
	return sqlDB.PingContext(ctx)
}

func (c *DbClientGorm) Select(ctx context.Context, table string, params models.FilterParams, out interface{}) error {
	exec := c.Db.WithContext(ctx).Table(table).Select(params.Select).Where(params.Filter).Scan(out)
	if exec.Error != nil {
		return exec.Error
	}
	if exec.RowsAffected == 0 {
		return errors.New(models.DB_ERROR_NOT_FOUND)
	}
	return nil
}

func (c *DbClientGorm) Get(ctx context.Context, query models.FilterParams, out interface{}) error {
	exec := c.Db.WithContext(ctx).Order(query.Orderings).Where(query.Filter).Limit(query.ValidLimit()).Offset(query.Offset).Find(out)

	if exec.Error != nil {
		return exec.Error
	}
	if exec.RowsAffected == 0 {
		return errors.New(models.DB_ERROR_NOT_FOUND)
	}
	return nil
}

func (c *DbClientGorm) GetTable(ctx context.Context, viewName string, params models.FilterParams, out interface{}) error {
	exec := c.Db.WithContext(ctx).Table(viewName).Order(params.Orderings).Where(params.Filter).Limit(params.ValidLimit()).Offset(params.Offset).Find(out)
	if exec.Error != nil {
		return exec.Error
	}
	if exec.RowsAffected == 0 {
		return errors.New(models.DB_ERROR_NOT_FOUND)
	}
	return nil
}

func (c *DbClientGorm) Create(ctx context.Context, input interface{}) error {
	if input == nil {
		return errors.New("creation data is nil")
	}
	return c.Db.WithContext(ctx).Create(input).Error
}

func (c *DbClientGorm) Update(ctx context.Context, params models.FilterParams, input interface{}) error {
	if input == nil {
		return errors.New("updated data is nil")
	}
	exec := c.Db.WithContext(ctx).Where(params.Filter).Updates(input)
	if exec.Error != nil {
		return exec.Error
	}
	if exec.RowsAffected == 0 {
		return errors.New(models.DB_ERROR_NOT_FOUND)
	}
	return nil
}

func (c *DbClientGorm) Upsert(ctx context.Context, params models.FilterParams, input interface{}) error {
	if input == nil {
		return errors.New("data is nil")
	}
	err := c.Update(ctx, params, input)
	if models.IsErrNotFound(err) {
		return c.Create(ctx, input)
	}
	return err
}

func (c *DbClientGorm) Delete(ctx context.Context, params models.FilterParams, input interface{}) error {
	if input == nil {
		return errors.New("delete data is nil")
	}
	exec := c.Db.WithContext(ctx).Where(params.Filter).Delete(input)
	if exec.Error != nil {
		return exec.Error
	}
	if exec.RowsAffected == 0 {
		return errors.New(models.DB_ERROR_NOT_FOUND)
	}
	return nil
}
