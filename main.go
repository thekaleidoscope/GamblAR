package main

import (
	"fmt"

	"github.com/GamblAR/atomic_elements"
	"github.com/pkg/errors"
)

func main() {
	game_handler := atomic_elements.Handeler{
		ChannelID: "townchannel",

		//OrgInfo
		OrgAdmins: map[string]string{
			"Org1": "Admin",
			"Org2": "Admin",
		},
		Organizations: []string{"Org1", "Org2"},
		ConfigFile:    "config.yaml",
	}

	fmt.Println("-------Initializing the game handler.--------")
	err := game_handler.Initializer()
	if err != nil {
		fmt.Println(errors.Wrap(err, "initialization of game handler failed"))

	} else {
		fmt.Println("-------Successful initializing the game handler.--------")
	}

	fmt.Println("------ Creating And Join Channel For each Organizations ----------")

	err = game_handler.CreateAndJoinChannel()
	if err != nil {
		fmt.Println(errors.Wrap(err, "CreateAndJoinChannel game handler failed"))

	} else {
		fmt.Println("-------Successful  Creating And Join Channel For each Organizations--------")
	}

}
