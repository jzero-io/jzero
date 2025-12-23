// Insert insert a new record into the database.
// Deprecated: use InsertV2 instead.
Insert(ctx context.Context, session sqlx.Session, data *{{.upperStartCamelObject}}) (sql.Result,error)
InsertV2(ctx context.Context, session sqlx.Session, data *{{.upperStartCamelObject}}) error