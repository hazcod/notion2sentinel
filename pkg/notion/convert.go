package notion

import (
	"encoding/json"
	"fmt"
	"time"
)

const (
	iso8601Format = "2006-01-02T15:04:05Z"
)

func (n *Notion) ConvertAuditLogsToMap(logs []AuditLog) ([]map[string]string, error) {
	var siemLogs []map[string]string

	for _, log := range logs {
		actorStr, err := json.Marshal(log.Actor)
		if err != nil {
			return nil, fmt.Errorf("could not marshal actor: %v", err)
		}

		changesStr, err := json.Marshal(log.Changes)
		if err != nil {
			return nil, fmt.Errorf("could not marshal changes: %v", err)
		}

		entry := map[string]string{
			"TimeGenerated": time.Unix(log.ServerTimestamp, 0).Format(iso8601Format),
			"Action":        log.ActionName,
			"Actor":         string(actorStr),
			"Changes":       string(changesStr),
			"ID":            log.ID,
		}

		siemLogs = append(siemLogs, entry)
	}

	return siemLogs, nil
}
