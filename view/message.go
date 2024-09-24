package view

import (
	"cmd-gram-cli/models"
	"fmt"
	"time"
)

func Messages(msg *models.MessageDTO, u *models.User) {
	if msg.UserID != u.ID {
		fmt.Printf("\t\t\t\t%s | %s\n", msg.Body, msg.Time.Format(time.RFC822))
	} else {
		fmt.Printf("%s | %s\n", msg.Time.Format(time.RFC822), msg.Body)
	}
}
