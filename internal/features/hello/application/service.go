package hello_app

type HelloUseCase interface {
	GetHello() string
}

type HelloService struct{}

func NewHelloService() *HelloService {
	return &HelloService{}
}

func (s *HelloService) GetHello() string {
	return "Hello, World!"
}
