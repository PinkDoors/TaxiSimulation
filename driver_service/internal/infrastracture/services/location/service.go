package location

import (
	"context"
	http_configs "driver_service/configs/http"
	"driver_service/internal/domain/models"
	"driver_service/internal/infrastracture/mappers"
	"driver_service/internal/infrastracture/services/location/dto"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"net/url"
	"strconv"
)

type Service struct {
	config *http_configs.Config
	logger *zap.Logger
	radius int
}

func NewService(
	config *http_configs.Config,
	logger *zap.Logger,
) *Service {
	return &Service{
		config: config,
		logger: logger,
		radius: 5000,
	}
}

func (s *Service) GetDriversWithLocations(ctx context.Context, latLngLiteral models.LatLngLiteral) ([]models.Driver, error) {
	endpoint := fmt.Sprintf("%s/location/v1/drivers", s.config.Url)

	latStringValue := strconv.FormatFloat(float64(latLngLiteral.Lat), 'f', -1, 32)
	lngStringValue := strconv.FormatFloat(float64(latLngLiteral.Lng), 'f', -1, 32)

	params := url.Values{}
	params.Add("lat", latStringValue)
	params.Add("lng", lngStringValue)
	params.Add("radius", strconv.Itoa(s.radius))

	fullURL := fmt.Sprintf("%s?%s", endpoint, params.Encode())

	resp, err := http.Get(fullURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		s.logger.Debug("Location-service returned " + strconv.Itoa(resp.StatusCode) + " code.")
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var driverDTOs []dto.Driver
	if err := json.NewDecoder(resp.Body).Decode(&driverDTOs); err != nil {
		return nil, err
	}

	return mappers.DriverWithLocationDtosToModels(driverDTOs)
}
