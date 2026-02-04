package apple

import (
	"io"
	"net/http"
)

// UploadAssets requests asset upload URLs from CloudKit.
func (c *cloudKit) UploadAssets(db CKDatabase, req *CKAssetsUploadRequest) (*CKAssetsUploadResponse, error) {
	var resp CKAssetsUploadResponse
	if err := c.doRequest("POST", db, "assets/upload", req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// UploadAssetData uploads binary data to the given upload URL.
func (c *cloudKit) UploadAssetData(uploadURL string, data io.Reader) error {
	req, err := http.NewRequest("PUT", uploadURL, data)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/octet-stream")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return &CKError{Code: CKErrorNetworkError, Reason: err.Error()}
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return &CKError{
			Code:   CKErrorUnknownError,
			Reason: "asset upload failed",
		}
	}

	return nil
}
