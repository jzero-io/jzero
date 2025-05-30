logic := {{ .OldService }}logic.New{{ .LogicTypeName }}(l.ctx, l.svcCtx)

// Create the adapter
adapter := &{{ .MethodName | FirstUpper }}ServerAdapter{
	stream,
}

marshal, err := proto.Marshal(in)
if err != nil {
	return err
}

var oldIn {{ .OldService }}pb.{{ .RequestTypeName }}
err = proto.Unmarshal(marshal, &oldIn)
if err != nil {
	return err
}

return logic.{{ .MethodName | FirstUpper }}(&oldIn, adapter)