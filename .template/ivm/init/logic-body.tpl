    logic := {{ .OldService }}logic.New{{ .LogicTypeName | FirstUpper }}(l.ctx, l.svcCtx)
    marshal, err := proto.Marshal(in)
    if err != nil {
    	return nil, err
    }
    var oldIn {{ .OldService }}pb.{{ .RequestTypeName }}
    err = proto.Unmarshal(marshal, &oldIn)
    if err != nil {
    	return nil, err
    }
    result, err := logic.{{ .MethodName }}(&oldIn)
    if err != nil {
    	return nil, err
    }
    marshal, err = proto.Marshal(result)
    if err != nil {
    	return nil, err
    }
    var newResp {{ .Service }}pb.{{ .ResponseTypeName }}
    err = proto.Unmarshal(marshal, &newResp)
    if err != nil {
    	return nil, err
    }
    return &newResp, nil