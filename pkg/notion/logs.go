package notion

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/hazcod/notion2sen/pkg/utils"
)

const (
	auditLogEndpoint = "https://www.notion.so/api/v3/searchAuditLogForOrganization"
)

type AuditLog struct {
	ActionName      string `json:"actionName"`
	ServerTimestamp int64  `json:"serverTimestamp"`
	OrganizationID  string `json:"organizationId"`
	Target          struct {
		Organization struct {
			Table string `json:"table"`
			ID    string `json:"id"`
		} `json:"organization"`
	} `json:"target,omitempty"`
	Changes any `json:"changes,omitempty"`
	Actor   struct {
		ActorID   string `json:"actorId"`
		ActorType string `json:"actorType"`
		Name      string `json:"name"`
		Email     string `json:"email"`
		Meta      struct {
			IPAddress        string `json:"ipAddress"`
			Platform         string `json:"platform"`
			City             string `json:"city"`
			CountryCode      string `json:"countryCode"`
			State            string `json:"state"`
			OrganizationRole string `json:"organizationRole"`
		} `json:"meta"`
	} `json:"actor"`
	Version int    `json:"version"`
	ID      string `json:"id"`
}
type auditLogResponse struct {
	Total          int        `json:"total"`
	Results        []AuditLog `json:"results"`
	HasMoreResults bool       `json:"hasMoreResults"`
	NextCursor     []any      `json:"nextCursor"`
}

func (n *Notion) GetAuditLogs(lookback time.Time) ([]AuditLog, error) {
	var auditLogs []AuditLog

	limit := 100
	var cursor []interface{}

	type auditLogRequest struct {
		Cursor  []interface{} `json:"cursor,omitempty"`
		Limit   int           `json:"limit"`
		Filters struct {
			TimeRange struct {
				Starting int64 `json:"starting"`
			} `json:"timeRange"`
		} `json:"filters"`
		Sort           string `json:"sort"`
		OrganisationID string `json:"organizationId"`
	}

	httpClient := utils.NewLogHttpClient(n.logger)

	for {
		logRequest := auditLogRequest{
			Limit: limit,
			Filters: struct {
				TimeRange struct {
					Starting int64 `json:"starting"`
				} `json:"timeRange"`
			}{
				TimeRange: struct {
					Starting int64 `json:"starting"`
				}{
					Starting: lookback.UnixMilli(),
				},
			},
			Sort:           "CreatedNewest",
			OrganisationID: n.orgID,
			Cursor:         cursor,
		}

		b, err := json.Marshal(&logRequest)
		if err != nil {
			return nil, fmt.Errorf("could not marshal audit log request: %v", err)
		}

		req, err := http.NewRequest(http.MethodPost, auditLogEndpoint, bytes.NewBuffer(b))
		if err != nil {
			return nil, fmt.Errorf("could not create http request: %v", err)
		}

		req.Header.Set("content-type", "application/json")
		req.Header.Set("authorization", fmt.Sprintf("Bearer %s", n.token))

		resp, err := httpClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("could not send http request: %v", err)
		}

		if resp.StatusCode > 399 {
			return nil, fmt.Errorf("could not get audit log: status=%v", resp.Status)
		}

		var auditLog auditLogResponse
		if err := json.NewDecoder(resp.Body).Decode(&auditLog); err != nil {
			return nil, fmt.Errorf("could not decode audit log response: %v", err)
		}

		auditLogs = append(auditLogs, auditLog.Results...)

		n.logger.Debugf("fetched page with total=%d results=%d cursor=%s", auditLog.Total, len(auditLog.Results), cursor)

		if !auditLog.HasMoreResults {
			n.logger.Debugf("no more results")
			break
		}

		cursor = auditLog.NextCursor
	}

	return auditLogs, nil
}
