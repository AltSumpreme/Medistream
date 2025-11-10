package workers

import (
	"github.com/hibiken/asynq"
)

func RegisterS3Workers(mux *asynq.ServeMux) {
	// mux.HandleFunc(string(JobTypeUploadReportPDF), handleUploadReportPDF)
}
