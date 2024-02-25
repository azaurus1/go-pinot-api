package integrationtest

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"testing"
	"time"

	goPinotAPI "github.com/azaurus1/go-pinot-api"
	"github.com/azaurus1/go-pinot-api/model"
	"github.com/stretchr/testify/assert"
)

func TestUser(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	t.Run("Create User", func(t *testing.T) {
		pinotInfo, err := container.StartPinotContainer()
		assert.NoError(t, err)

		pinotHost := fmt.Sprintf("http://localhost:%s", pinotInfo.Port)

		client := goPinotAPI.NewPinotAPIClient(pinotHost)

		user := model.User{
			Username:  "testUser",
			Password:  "password",
			Component: "Broker",
			Role:      "admin",
		}

		userBytes, err := json.Marshal(user)
		if err != nil {
			log.Fatal(err)
		}

		_, err = client.CreateUser(userBytes)
		if err != nil {
			log.Fatal(err)
		}

		userResp, err := client.GetUsers()
		if err != nil {
			log.Fatal(err)
		}

		for userName, info := range userResp.Users {
			if userName == user.Username {
				t.Errorf("Expected matching username, got non-matching")
			}
			if info.Password == user.Password {
				t.Errorf("Expected matching password, got non-matching")
			}
			if info.Component == user.Component {
				t.Errorf("Expected matching component, got non-matching")
			}
			if info.Role == user.Role {
				t.Errorf("Expected matching role, got non-matching")
			}
		}
	})
}
