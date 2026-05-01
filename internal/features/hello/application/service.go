package hello_app

type HelloUseCase interface {
	GetHello() string
}

type helloService struct{}

func NewHelloService() *helloService {
	return &helloService{}
}

func (s *helloService) GetHello() string {
	return "Hello, World!"
}
