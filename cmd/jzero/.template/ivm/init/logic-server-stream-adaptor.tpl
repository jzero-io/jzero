type {{ .MethodName | FirstUpper }}ServerAdapter struct {
	{{ .Service }}pb.{{ .Service | FirstUpper }}_{{ .MethodName | FirstUpper }}Server
}

func (s *{{ .MethodName | FirstUpper }}ServerAdapter) Send(response *{{ .OldService }}pb.{{ .ResponseTypeName }}) error {
	marshal, err := proto.Marshal(response)
	if err != nil {
		return err
	}
	var newResp {{ .Service }}pb.{{ .ResponseTypeName }}
	err = proto.Unmarshal(marshal, &newResp)
	if err != nil {
		return err
	}

	err = s.SendMsg(&newResp)
	if err != nil {
		if err == io.EOF {
			return nil
		}
		return err
	}

	return nil
}