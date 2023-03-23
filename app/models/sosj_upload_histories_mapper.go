package models

func (result *GetSosjUploadHistoryResponse) GetSosjUploadHistoryResponseMap(request *UploadHistory) {
	result.ID = request.ID
	result.RequestId = request.RequestId
	result.BulkCode = request.BulkCode
	result.FileName = request.FileName
	result.FilePath = request.FilePath
	result.AgentId = request.AgentId
	result.AgentName = request.AgentName
	result.UploadedBy = request.UploadedBy
	result.UploadedByName = request.UploadedByName
	result.UploadedByEmail = request.UploadedByEmail
	result.UpdatedBy = request.UpdatedBy
	result.UpdatedByName = request.UpdatedByName
	result.UpdatedByEmail = request.UpdatedByEmail
	result.Status = request.Status
	result.TotalRows = request.TotalRows
	result.CreatedAt = request.CreatedAt
	result.UpdatedAt = request.UpdatedAt
}
