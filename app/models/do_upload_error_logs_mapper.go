package models

import (
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (result *DoUploadErrorLog) DoUploadErrorLogsMap(line int, doUploadHistoryId, requestId, bulkCode string, errors []string, now *time.Time) {
	result.RequestId = requestId
	result.DoUploadHistoryId, _ = primitive.ObjectIDFromHex(doUploadHistoryId)
	result.BulkCode = bulkCode
	result.ErrorRowLine = int64(line)
	result.ErrorMessage = strings.Join(errors, ";")
	result.CreatedAt = now
	result.UpdatedAt = now
}
