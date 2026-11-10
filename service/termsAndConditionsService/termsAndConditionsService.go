package termsAndConditionsService

import (
	"Quote/db/daoPostgres/daos/termsAndConditionsDao"
	"Quote/service/fileService"
	"context"
)

// GetPdf returns the terms and conditions PDF, or nil if none has been uploaded.
func GetPdf(ctx context.Context) ([]byte, error) {
	terms, found, err := termsAndConditionsDao.GetInstance().Get(ctx)
	if err != nil || !found {
		return nil, err
	}

	data, found, err := fileService.GetSingleFile(ctx, terms.FileID)
	if err != nil || !found {
		return nil, err
	}

	return data, nil
}
