package svc

type Middleware struct{}

func (svcCtx *ServiceContext) NewMiddleware() Middleware {
	return Middleware{}
}
