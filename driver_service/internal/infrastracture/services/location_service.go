package services

type LocationService struct {
	// Ваши поля для настройки клиента
}

// GetData реализует метод интерфейса ExternalServiceClient.
func (c *LocationService) GetData() (string, error) {
	// Ваша логика для получения данных от внешнего микросервиса
	return "data from external service", nil
}
