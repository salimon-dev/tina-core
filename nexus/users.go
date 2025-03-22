package nexus

import (
	"encoding/json"
	"fmt"
	"salimon/tina-core/types"

	"github.com/google/uuid"
)

func FetchUserData(userId uuid.UUID) (*types.UserData, error) {
	url := "/entity/user/" + userId.String()
	response, err := sendHttpRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	var user types.UserData
	err = json.Unmarshal(response, &user)
	return &user, err
}
