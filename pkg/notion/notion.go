package notion

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

type Notion struct {
	logger *logrus.Logger
	token  string
	orgID  string
}

func New(logger *logrus.Logger, token, organisationID string) (*Notion, error) {
	if logger == nil {
		logger = logrus.New()
	}

	if token == "" {
		return nil, fmt.Errorf("no notion token provided")
	}

	if organisationID == "" {
		return nil, fmt.Errorf("no notion organisation ID provided")
	}

	return &Notion{
		logger: logger,
		token:  token,
		orgID:  organisationID,
	}, nil
}
