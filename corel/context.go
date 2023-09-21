package corel

import (
	"context"
)

func NewContext(corelid *CoRelationId) context.Context {
	return context.WithValue(context.Background(), CtxCorelLocator, corelid)
}

/*
-- if len(ids)>0 then a new corel will be created and this will be set in cotenxt.WithValue with the createdCorelId attached
-- else if the passed context have a corel,
*The return context is context.WithValue and in th corel it has the child corel of that of passed corel context
*/
func NewChildContext(ctx context.Context, ids ...string) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	var corelid *CoRelationId
	if len(ids) > 0 {
		corelid = NewCorelId(ids...)
	} else if corelid1, ok := ctx.Value(CtxCorelLocator).(*CoRelationId); ok {
		corelid = corelid1.Child()
	} else {
		corelid = NewCorelId()
	}
	return context.WithValue(ctx, CtxCorelLocator, corelid)
}

func NewOrphenContext(ids ...string) context.Context {
	return NewChildContext(context.Background(), ids...)
}
