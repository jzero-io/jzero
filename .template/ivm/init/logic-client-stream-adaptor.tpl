type {{ .MethodName | FirstUpper }}ServerAdapter struct {
	{{ .Service }}pb.{{ .Service | FirstUpper }}_{{ .MethodName | FirstUpper }}Server
}

func (s *{{ .MethodName | FirstUpper }}ServerAdapter) SendAndClose(response *{{ .OldService }}pb.{{ .ResponseTypeName }}) error {
	marshal, err := proto.Marshal(response)
	if err != nil {
		return err
	}

	var newResp {{ .Service }}pb.{{ .ResponseTypeName }}
	err = proto.Unmarshal(marshal, &newResp)
	if err != nil {
		return err
	}
	return s.SendMsg(&newResp)
}

func (s *{{ .MethodName | FirstUpper }}ServerAdapter) Recv() (*{{ .OldService }}pb.{{ .RequestTypeName }}, error) {
	for {
		newIn, err := s.{{ .Service | FirstUpper }}_{{ .MethodName | FirstUpper }}Server.Recv()
		if err == io.EOF {
			return nil, io.EOF
		}
		if err != nil {
			return nil, err
		}

		marshal, err := proto.Marshal(newIn)
		if err != nil {
			return nil, err
		}
		var oldIn {{ .OldService }}pb.{{ .RequestTypeName }}
		err = proto.Unmarshal(marshal, &oldIn)
		if err != nil {
			return nil, err
		}

		return &oldIn, nil
	}
}