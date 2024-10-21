package auth

import (
	"datafetcher/storage/cloud"
	"fmt"
	"log"
)

func CreateOAuth2(cloudID int) OAuth2 {
	switch cloudID {
	case cloud.Google:
		{
			return OAuth2Google{}
		}
	default:
		{
			mess := fmt.Sprintf("The cloud ID provided does not match any known supported cloud: %d", cloudID)
			log.Println(mess)
			return nil
		}
	}
}
