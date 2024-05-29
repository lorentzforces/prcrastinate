package github

import (
	"context"
	"fmt"
	"time"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

const defaultTimeout = time.Second * 5

type GhClient struct {
	extClient githubv4.Client
}

func GetClient(token string) *GhClient {
	authSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	// provided context determines the valid lifetime for the client (forever)
	httpClient := oauth2.NewClient(context.Background(), authSource)
	// TODO: determine if we need to switch endpoints for GH enterprise cloud

	return &GhClient {
		extClient: *githubv4.NewClient(httpClient),
	}
}

// reusable type to extract a username from a Github User structure
type UserLogin struct {
	Login string
}

type UserQuery struct {
	Viewer UserLogin
}

func (u *UserQuery) mapToUser() *User {
	return &User {
		Name: u.Viewer.Login,
	}
}

type User struct {
	Name string
}

func (client *GhClient) FetchUser() (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 5)
	defer cancel()

	var query UserQuery

	err := client.extClient.Query(ctx, &query, nil)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch user from Github: %w", err)
	}

	return query.mapToUser(), nil
}

type PrQuery struct {
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
					Author UserLogin
					CreatedAt time.Time
					UpdatedAt time.Time
					Repository struct {
						Owner UserLogin
						Name string
					}
					Reviews struct {
						Nodes []struct {
							Id string
							Url string
							PublishedAt time.Time
							UpdatedAt time.Time
							State string
							Author UserLogin
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

func (p *PrQuery) mapToReviewData() *ReviewData {
	return &ReviewData {
		PrCount: p.Search.IssueCount,
	}
}

type ReviewData struct {
	PrCount int32
}

func (client *GhClient) FetchPrData(username string) (*ReviewData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 5)
	defer cancel()

	var query PrQuery
	searchStr := fmt.Sprintf("type:pr state:closed author:%s", username)
	queryArgs := map[string]interface{} {
		"search_str": githubv4.String(searchStr),
		"curr_cursor": githubv4.String(""),
	}
	err := client.extClient.Query(ctx, &query, queryArgs)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch review data from Github: %w", err)
	}

	return query.mapToReviewData(), nil
}
