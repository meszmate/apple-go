package apple

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestQueryRecords(t *testing.T) {
	mockedClient := new(MockedCKHTTPClient)
	ck := newTestCloudKit(mockedClient)

	resp := CKQueryResponse{
		Records: []CKRecord{
			{RecordName: "rec-1", RecordType: "TestType"},
			{RecordName: "rec-2", RecordType: "TestType"},
		},
	}
	respBytes, _ := json.Marshal(resp)

	mockedClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		return req.Method == "POST" && req.URL.Path == "/database/1/iCloud.com.example.app/development/public/records/query"
	})).Return(
		&http.Response{StatusCode: 200, Body: io.NopCloser(bytes.NewReader(respBytes))},
		nil,
	)

	result, err := ck.QueryRecords(CKDatabasePublic, &CKQueryRequest{
		Query: CKQuery{RecordType: "TestType"},
	})
	assert.NoError(t, err)
	assert.Len(t, result.Records, 2)
	assert.Equal(t, "rec-1", result.Records[0].RecordName)
}

func TestQueryRecordsError(t *testing.T) {
	mockedClient := new(MockedCKHTTPClient)
	ck := newTestCloudKit(mockedClient)

	errResp := CKErrorResponse{ServerErrorCode: CKErrorNotFound, Reason: "not found"}
	errBytes, _ := json.Marshal(errResp)

	mockedClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		return true
	})).Return(
		&http.Response{StatusCode: 404, Body: io.NopCloser(bytes.NewReader(errBytes))},
		nil,
	)

	result, err := ck.QueryRecords(CKDatabasePublic, &CKQueryRequest{
		Query: CKQuery{RecordType: "TestType"},
	})
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestModifyRecords(t *testing.T) {
	mockedClient := new(MockedCKHTTPClient)
	ck := newTestCloudKit(mockedClient)

	resp := CKRecordsModifyResponse{
		Records: []CKRecord{{RecordName: "new-rec", RecordType: "TestType"}},
	}
	respBytes, _ := json.Marshal(resp)

	mockedClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		return req.Method == "POST" && req.URL.Path == "/database/1/iCloud.com.example.app/development/public/records/modify"
	})).Return(
		&http.Response{StatusCode: 200, Body: io.NopCloser(bytes.NewReader(respBytes))},
		nil,
	)

	result, err := ck.ModifyRecords(CKDatabasePublic, &CKRecordsModifyRequest{
		Operations: []CKRecordOperation{
			{
				OperationType: CKOperationCreate,
				Record: CKRecord{
					RecordType: "TestType",
					Fields: map[string]*CKField{
						"name": {Value: "test", Type: "STRING"},
					},
				},
			},
		},
	})
	assert.NoError(t, err)
	assert.Len(t, result.Records, 1)
	assert.Equal(t, "new-rec", result.Records[0].RecordName)
}

func TestModifyRecordsError(t *testing.T) {
	mockedClient := new(MockedCKHTTPClient)
	ck := newTestCloudKit(mockedClient)

	errResp := CKErrorResponse{ServerErrorCode: CKErrorConflict, Reason: "conflict"}
	errBytes, _ := json.Marshal(errResp)

	mockedClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		return true
	})).Return(
		&http.Response{StatusCode: 409, Body: io.NopCloser(bytes.NewReader(errBytes))},
		nil,
	)

	result, err := ck.ModifyRecords(CKDatabasePublic, &CKRecordsModifyRequest{
		Operations: []CKRecordOperation{
			{OperationType: CKOperationUpdate, Record: CKRecord{RecordName: "rec"}},
		},
	})
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestLookupRecords(t *testing.T) {
	mockedClient := new(MockedCKHTTPClient)
	ck := newTestCloudKit(mockedClient)

	resp := CKRecordsLookupResponse{
		Records: []CKRecord{{RecordName: "rec-1"}},
	}
	respBytes, _ := json.Marshal(resp)

	mockedClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		return req.URL.Path == "/database/1/iCloud.com.example.app/development/public/records/lookup"
	})).Return(
		&http.Response{StatusCode: 200, Body: io.NopCloser(bytes.NewReader(respBytes))},
		nil,
	)

	result, err := ck.LookupRecords(CKDatabasePublic, &CKRecordsLookupRequest{
		Records: []CKRecord{{RecordName: "rec-1"}},
	})
	assert.NoError(t, err)
	assert.Len(t, result.Records, 1)
}

func TestRecordChanges(t *testing.T) {
	mockedClient := new(MockedCKHTTPClient)
	ck := newTestCloudKit(mockedClient)

	resp := CKRecordChangesResponse{
		Records:   []CKRecord{{RecordName: "changed-rec"}},
		SyncToken: "new-sync-token",
	}
	respBytes, _ := json.Marshal(resp)

	mockedClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		return req.URL.Path == "/database/1/iCloud.com.example.app/development/public/records/changes"
	})).Return(
		&http.Response{StatusCode: 200, Body: io.NopCloser(bytes.NewReader(respBytes))},
		nil,
	)

	result, err := ck.RecordChanges(CKDatabasePublic, &CKRecordChangesRequest{
		ZoneID:    CKZoneID{ZoneName: "_defaultZone"},
		SyncToken: "old-sync-token",
	})
	assert.NoError(t, err)
	assert.Len(t, result.Records, 1)
	assert.Equal(t, "new-sync-token", result.SyncToken)
}
