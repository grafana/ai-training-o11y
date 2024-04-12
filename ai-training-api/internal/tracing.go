package db

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

// Defining a type and iota value creates a namespaced key in a context.Context.
type contextKey int

const (
	spanKey contextKey = iota
)

type tracer struct {
}

func (t *tracer) Name() string {
	return "gorm:tracing"
}

//nolint:errcheck // Error checking would make this function much harder to understand.
func (t *tracer) Initialize(db *gorm.DB) error {
	db.Callback().Create().Before("*").Register("trace:create_before", before("gorm_create"))
	db.Callback().Create().After("*").Register("trace:create_after", after)

	db.Callback().Query().Before("*").Register("trace:query_before", before("gorm_query"))
	db.Callback().Query().After("*").Register("trace:query_after", after)

	db.Callback().Update().Before("*").Register("trace:update_before", before("gorm_update"))
	db.Callback().Update().After("*").Register("trace:update_after", after)

	db.Callback().Delete().Before("*").Register("trace:delete_before", before("gorm_delete"))
	db.Callback().Delete().After("*").Register("trace:delete_after", after)

	db.Callback().Row().Before("*").Register("trace:row_before", before("gorm_row"))
	db.Callback().Row().After("*").Register("trace:row_after", after)

	db.Callback().Raw().Before("*").Register("trace:raw_before", before("gorm_raw"))
	db.Callback().Raw().After("*").Register("trace:raw_after", after)
	return nil
}

// before creates a callback that starts a span and adds it to the context.
func before(operationName string) func(*gorm.DB) {
	return func(db *gorm.DB) {
		ctx := db.Statement.Context
		ctx, span := otel.Tracer("db").Start(ctx, operationName)
		ctx = context.WithValue(ctx, spanKey, span)
		db.Statement.Context = ctx
	}
}

// after finishes the span attached to the db's context.
func after(db *gorm.DB) {
	span, ok := db.Statement.Context.Value(spanKey).(trace.Span)
	if !ok {
		return
	}
	defer span.End()

	if db.Error != nil {
		span.RecordError(db.Error)
	}
}
