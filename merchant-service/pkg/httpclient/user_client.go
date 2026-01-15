package httpclient

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"micro-warehouse/merchant-service/configs"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2/log"
)

type UserClientInterface interface {
	GetUserByID(ctx context.Context, userID uint) (*UserResponse, error)
}

type UserClient struct {
	urlUserService string
	httpClient     *http.Client
}

type UserResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
	Photo string `json:"photo"`
}

type UserServiceResponse struct {
	Message string       `json:"message"`
	Data    UserResponse `json:"data"`
	Error   string       `json:"error,omitempty"`
}

func NewUserClient(cfg configs.Config) UserClientInterface {
	return &UserClient{
		urlUserService: cfg.App.UrlUserService,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GetUserByID implements [UserClientInterface].
func (u *UserClient) GetUserByID(ctx context.Context, userID uint) (*UserResponse, error) {
	url := fmt.Sprintf("%s/api/v1/users/%d", u.urlUserService, userID)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		log.Errorf("[UserClient] GetUserByID - 1: %v", err)
		return nil, err
	}

	resp, err := u.httpClient.Do(req)
	if err != nil {
		log.Errorf("[UserClient] GetUserByID - 2: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Errorf("[UserClient] GetUserByID - 3: received non-200 response code: %d", resp.StatusCode)
		return nil, fmt.Errorf("failed to get user: status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("[UserClient] GetUserByID - 4: %v", err)
		return nil, err
	}

	var userResp UserServiceResponse
	err = json.Unmarshal(body, &userResp)
	if err != nil {
		log.Errorf("[UserClient] GetUserByID - 5: %v", err)
		return nil, err
	}

	return &userResp.Data, nil
}
