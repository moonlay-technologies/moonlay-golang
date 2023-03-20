package models

import (
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (result *SosjUploadErrorLog) SosjUploadErrorLogsMap(request map[string]string, sosjUploadHistoryId, requestId, agentName string, line int, errors []string, now *time.Time) {
	result.RequestId = requestId
	result.SosjUploadHistoryId, _ = primitive.ObjectIDFromHex(sosjUploadHistoryId)
	result.BulkCode = "SOSJ-" + request["IDDistributor"] + "-" + fmt.Sprint(now.Unix())
	result.ErrorRowLine = int64(line)
	result.ErrorMessage = strings.Join(errors, ";")
	result.CreatedAt = now
	result.UpdatedAt = now
}

func (result *RowDateSosjUploadErrorLog) RowDateSosjUploadErrorLogMap()
