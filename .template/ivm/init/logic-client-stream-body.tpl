logic := {{ .OldService }}logic.New{{ .LogicTypeName }}(l.ctx, l.svcCtx)

// Create the adapter
adapter := &{{ .MethodName | FirstUpper }}ServerAdapter{
    stream,
}

return logic.{{ .MethodName | FirstUpper }}(adapter)