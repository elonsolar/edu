package repo

import (
	"context"
	"database/sql/driver"
	"edu/internal/conf"
	"edu/internal/domain"
	"edu/internal/domain/model"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	"github.com/google/wire"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewGormDB, NewBaseRepo, NewTxImp, NewUserRepo, NewTenantRepo, NewTenantPermissionRepo, NewTeacherRepo, NewStudentRepo, NewRoomRepo, NewSubjectRepo, NewCustomerRepo, NewLessonHistoryRepo, NewCoursePlanRepo, NewCoursePlanDetailRepo, NewCoursePlanStudentRepo, NewDailyLessonRepo, NewDailyLessonStudentRepo, NewMetaRepo, NewPermissionRepo, NewRoleRepo, NewRolePermissionRepo, NewTaskRepo, NewSkuRepo, NewCombineSkuRepo, NewCombineSkuItemRepo)

type Entity struct {
	gorm.Model
	Version int32
}

func (e *Entity) setId(id int32) {
	e.ID = uint(id)
}

type BaseRepo struct {
	*gorm.DB
	Log *log.Helper
}

func NewBaseRepo(c *conf.Data, db *gorm.DB, logger log.Logger) (*BaseRepo, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	return &BaseRepo{DB: db, Log: log.NewHelper(logger)}, cleanup, nil
}

// 对于 普通的db 方法，如果ctx 中有事务，则沿用事务对象
// ctx 没有事务，则新建 gorm.DB
func (d *BaseRepo) GetDBFromContext(ctx context.Context) *gorm.DB {

	ctxVal := ctx.Value(ctxTransactionKey{})

	if ctxVal == nil {
		return d.DB.WithContext(ctx)
	} else {
		return ctxVal.(*gorm.DB)
	}
}

func NewGormDB(c *conf.Data) *gorm.DB {
	dsn := "root:123456@tcp(127.0.0.1:3306)/we_friends?charset=utf8mb4&parseTime=True&loc=Local"
	// dsn := c.Database
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
	if err != nil {
		panic(err)
	}
	// 注册插件
	db.Use(&OptimisticLock{})

	return db
}

// design 1:  Repo.Begin(ctx) :a
//			  xRepo.DoA(ctx) :b
//			  yRepo.DoB(ctx) :c
//			  Repo.Commit(ctx) :d

// a 必须返回 ctx,因为ctx 存值后必须产生一个新的context

// design 2:
//	 		此种方式可以支持事物嵌套， 可以判断 是否是事务开起方（事务开起方才可以 commit ）
//          xRepo.Transaction(ctx ,func(cxt)error{
//
//
//
//          })

type ctxTransactionKey struct{}
type TxImp struct {
	d   *gorm.DB
	log log.Logger
}

// Transaction implements domain.Tx.
func (t *TxImp) Transaction(ctx context.Context, fn func(txctx context.Context) error) error {
	ctxVal := ctx.Value(ctxTransactionKey{})

	var tx *gorm.DB

	var isOpener = false

	if ctxVal == nil {
		tx = t.d.WithContext(ctx).Begin()

		ctx = context.WithValue(ctx, ctxTransactionKey{}, tx)
		isOpener = true
	} else {
		tx = ctxVal.(*gorm.DB)
	}

	err := fn(ctx)
	if err != nil {
		return tx.Rollback().AddError(err)
	}

	if isOpener {
		return tx.Commit().Error
	}
	return nil
}

func NewTxImp(d *gorm.DB, log log.Logger) domain.Tx {

	return &TxImp{d, log}
}

func query2Clause(expr *domain.Expression) clause.Expression {

	var defalutExpr = clause.Expr{
		SQL:                "1=1",
		Vars:               []interface{}{},
		WithoutParentheses: false,
	}

	if expr == nil {
		return defalutExpr
	}
	if expr.IsLogic {
		if expr.Op == "and" {
			return clause.And(parseExpr(expr.SubExpr)...)
		} else {
			return clause.Or(parseExpr(expr.SubExpr)...)
		}
	}
	return defalutExpr
}

func parseExpr(exprs []*domain.Expression) []clause.Expression {

	var result = []clause.Expression{}
	for _, expr := range exprs {

		if expr.IsLogic {
			if expr.Op == "and" {
				return []clause.Expression{clause.And(parseExpr(expr.SubExpr)...)}
			} else {
				return []clause.Expression{clause.Or(parseExpr(expr.SubExpr)...)}
			}
		} else {
			if expr.Value != nil {
				switch expr.Op {
				case "=":
					result = append(result, clause.Eq{
						Column: expr.Column,
						Value:  expr.Value,
					})
				case "like":
					result = append(result, clause.Like{
						Column: expr.Column,
						Value:  "%" + fmt.Sprintf("%v", expr.Value) + "%",
					})
				case "in":

					if values, ok := expr.Value.([]interface{}); ok {
						result = append(result, clause.IN{
							Column: expr.Column,
							Values: values,
						})
					}

				}

			}

		}
	}
	return result
}

func Map[S any, T any](source []S, fn func(S) T) []T {

	var result []T
	for _, s := range source {
		result = append(result, fn(s))
	}
	return result
}

func intVal(name string, params map[string]string, defalut int) int {

	val, ok := params[name]
	if !ok {
		return defalut
	}
	intVal, err := strconv.Atoi(val)
	if err != nil {
		return defalut
	}
	return intVal
}

// scope
func Paginate(page *domain.Page) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page == nil {
			page = &domain.Page{}
		}

		var pageNo = page.PageNo
		var pageSize = page.PageSize

		if pageNo < 1 {
			pageNo = 1
		}

		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (pageNo - 1) * pageSize
		return db.Offset(int(offset)).Limit(int(pageSize))
	}

}

type OptimisticLock struct {
}

func (ol *OptimisticLock) Name() string {
	return "Optimistic"
}

func (ol *OptimisticLock) Initialize(db *gorm.DB) error {

	db.Callback().Update().Before("gorm:update").Register("optimistic:before_update", func(d *gorm.DB) {
		if d.Statement.Schema == nil {
			return
		}

		destVal := reflect.Indirect(reflect.ValueOf(d.Statement.Dest))
		switch destVal.Kind() {
		case reflect.Struct:
			field := d.Statement.Schema.LookUpField("version")
			// Get value from field
			if fieldValue, isZero := field.ValueOf(d.Statement.Context, destVal); !isZero {
				if version, ok := fieldValue.(int32); ok {
					d.Where("version=?", version)
					// Set value to field
					err := field.Set(d.Statement.Context, destVal, version+1)
					if err != nil {
						d.Logger.Error(d.Statement.Context, "set version ", err)
					}
				}
			}

		}
	})

	db.Callback().Create().Before("gorm:create").Register("optimistic:before_create", func(d *gorm.DB) {
		tenantId := d.Statement.Schema.LookUpField("TenantId")

		if tenantId == nil {
			return
		}

		reflectValue := reflect.Indirect(reflect.ValueOf(d.Statement.Dest))
		for reflectValue.Kind() == reflect.Ptr || reflectValue.Kind() == reflect.Interface {
			reflectValue = reflect.Indirect(reflectValue)
		}
		claims, ok := jwt.FromContext(d.Statement.Context)
		if !ok {
			return
		}
		c, ok := (claims).(*model.Claims)
		if !ok {
			d.Logger.Error(d.Statement.Context, "get claim err")
			return
		}
		if c.IsSuperAdmin {
			return
		}

		switch reflectValue.Kind() {
		case reflect.Slice, reflect.Array:
			for i := 0; i < reflectValue.Len(); i++ {
				var elementValue = reflectValue.Index(i)

				for elementValue.Kind() == reflect.Ptr || elementValue.Kind() == reflect.Interface {
					elementValue = reflect.Indirect(elementValue)
				}
				tenantId.Set(d.Statement.Context, elementValue, c.TenantId)
			}
		case reflect.Struct:
			tenantId.Set(d.Statement.Context, reflectValue, c.TenantId)
		}
	})

	db.Callback().Query().Before("gorm:query").Register("optimistic:before_query", func(d *gorm.DB) {

		tenantId := d.Statement.Schema.LookUpField("TenantId")

		if tenantId == nil {
			return
		}
		claims, ok := jwt.FromContext(d.Statement.Context)
		if !ok {
			return
		}
		c, ok := (claims).(*model.Claims)
		if !ok {
			d.Logger.Error(d.Statement.Context, "get claim err")
			return
		}
		if c.IsSuperAdmin {
			return
		}

		d.Statement.AddClause(clause.Where{
			Exprs: []clause.Expression{clause.Eq{
				Column: "tenant_id",
				Value:  c.TenantId,
			}},
		})

	})

	return nil

}

// types
type SliceInt32 []int32

func (j *SliceInt32) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSON value:", value))
	}

	result := SliceInt32{}
	err := json.Unmarshal(bytes, &result)
	*j = result
	return err
}

// Value return json value, implement driver.Valuer interface
func (j SliceInt32) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return json.Marshal(j)
}

type Slice[T any] []T

func (j *Slice[T]) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSON value:", value))
	}

	result := Slice[T]{}
	err := json.Unmarshal(bytes, &result)
	*j = result
	return err
}

// Value return json value, implement driver.Valuer interface
func (j Slice[T]) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return json.Marshal(j)
}
