package service

type Service interface {
	SimpleData(username string) string
	AdminData(username string) string

	// HealthCheck check service health status 健康检测
	HealthCheck() bool
}

type CommonService struct {
}

func (c *CommonService) SimpleData(username string) string {
	return "hello " + username + " ,simple data, with simple authority"
}

func (c *CommonService) AdminData(username string) string {
	return "hello " + username + " ,admin data, with admin authority"
}

func (c *CommonService) HealthCheck() bool {
	return true
}

func NewCommonService() *CommonService {
	return &CommonService{}
}
