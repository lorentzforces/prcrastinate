package command

import (
	"fmt"
	"prcrastinate/internal/github"
	"prcrastinate/internal/platform"
)

type Pull struct{}

func (cmd Pull) Name() string {
	return "pull"
}

func (cmd Pull) ShortDescr() string {
	return "Pull latest PR information from Github"
}

func (cmd Pull) usageStr() string {
	return "PLACEHOLDER usage for [pull]"
}

type PullArgs struct {
	BaseArgs
}

func (cmd Pull) Run(args []string) {
	parsedArgs := new(PullArgs)
	baseFlags := InitBaseFlags(&parsedArgs.BaseArgs)
	baseFlags.Parse(args)
	config := platform.ReadConfigFromPath(parsedArgs.ConfigPath)
	client := github.GetClient(config.Token)

	user, err := client.FetchUser()
	if err != nil {
		platform.FailOut(err.Error())
	}
	fmt.Printf("==DEBUG== SUCCESS! username: %s\n", user.Name)

	prData, err:= client.FetchPrData(user.Name)
	if err != nil {
		platform.FailOut(err.Error())
	}
	fmt.Printf("==DEBUG== SUCCESS! PR count: %d\n", prData.PrCount)

	// TODO: refresh local db
	// TODO: print stats?
}
