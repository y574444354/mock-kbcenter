package thirdPlatform

import (
	"github.com/zgsm/mock-kbcenter/pkg/httpclient"
)

var TypeIssueManager = "issueManager"

type IssueManagerService struct {
	*Service
}

func NewIssueManagerService() (*IssueManagerService, error) {
	clientConfig, err := GetServiceConfig(TypeIssueManager)
	if err != nil {
		return nil, err
	}

	client, err := httpclient.NewClient(clientConfig)
	if err != nil {
		return nil, err
	}

	service := &IssueManagerService{
		Service: &Service{
			client: client,
		},
	}

	return service, nil
}
