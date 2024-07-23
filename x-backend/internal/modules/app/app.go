package app

type controller struct {
	service service
}

// Init создает сервис приложения
func Init() *controller {
	service := InitService()

	return &controller{
		service: *service,
	}
}
