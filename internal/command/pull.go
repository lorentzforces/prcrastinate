package command

import (
	"context"
	"fmt"
	"prcrastinate/internal/platform"
	"time"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
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

	authSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: config.Token})
	// provided context determines the valid lifetime for the client
	httpClient := oauth2.NewClient(context.Background(), authSource)
	// TODO: determine if we need to switch endpoints for GH enterprise cloud
	ghClient := githubv4.NewClient(httpClient)

	var query struct {
		Viewer struct {
			Login     githubv4.String
			CreatedAt githubv4.DateTime
		}
	}

	userFetchCtx, userFetchCancel := context.WithTimeout(context.Background(), time.Second * 5)
	err := ghClient.Query(userFetchCtx, &query, nil)
	if err != nil {
		userFetchCancel()
		platform.FailOut(fmt.Sprintf("Failed to query Github: %s", err.Error()))
	}

	fmt.Printf("==DEBUG== SUCCESS! username: %s\n", query.Viewer.Login)
	userFetchCancel()

	// TODO: fetch relevant PRs
	// TODO: make this query type declaration a little more sane
	var prQuery struct {
		Search struct {
			PageInfo struct {
				StartCursor string
				HasNextPage bool
			}
			IssueCount int32
			Edges []struct{
				Node struct {
					PullRequest struct {
						Number int32
						Title string
						Author struct { // TODO: probably define this as a reusable type
							Login string
						}
						CreatedAt time.Time
						UpdatedAt time.Time
						Repository struct {
							Owner struct {
								Login string
							}
							Name string
						}
						Reviews struct {
							Nodes []struct {
								Id string
								Url string
								PublishedAt time.Time
								UpdatedAt time.Time
								State string
								Author struct {
									Login string
								}
								ViewerDidAuthor bool
								BodyText string
								Comments struct {
									TotalCount int
								}
							}
						} `graphql:"reviews(first: 50)"`
					} `graphql:"... on PullRequest"`
				}
			}
		} `graphql:"search(query: $search_str, type: ISSUE, first: 50, after: $curr_cursor)"`
	}

	prArgs := map[string]interface{} {
		// TODO: parameterize the username here
		"search_str": githubv4.String("type:pr state:closed author:lorentzforces"),
		"curr_cursor": githubv4.String(""),
	}

	reqCtx, cancelFunc := context.WithTimeout(context.Background(), time.Second * 5)
	defer cancelFunc()
	err = ghClient.Query(reqCtx, &prQuery, prArgs)
	if err != nil {
		platform.FailOut(err.Error())
	}

	fmt.Printf("==DEBUG== SUCCESS! PR count: %d\n", prQuery.Search.IssueCount)
	fmt.Printf("%+v\n", prQuery)

	// TODO: refresh local db
	// TODO: print stats?
}
