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

func TestListZones(t *testing.T) {
	mockedClient := new(MockedCKHTTPClient)
	ck := newTestCloudKit(mockedClient)

	resp := CKZonesResponse{
		Zones: []CKZone{
			{ZoneID: CKZoneID{ZoneName: "_defaultZone"}},
			{ZoneID: CKZoneID{ZoneName: "CustomZone"}},
		},
	}
	respBytes, _ := json.Marshal(resp)

	mockedClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		return req.URL.Path == "/database/1/iCloud.com.example.app/development/private/zones/list"
	})).Return(
		&http.Response{StatusCode: 200, Body: io.NopCloser(bytes.NewReader(respBytes))},
		nil,
	)

	result, err := ck.ListZones(CKDatabasePrivate)
	assert.NoError(t, err)
	assert.Len(t, result.Zones, 2)
	assert.Equal(t, "_defaultZone", result.Zones[0].ZoneID.ZoneName)
}

func TestListZonesError(t *testing.T) {
	mockedClient := new(MockedCKHTTPClient)
	ck := newTestCloudKit(mockedClient)

	errResp := CKErrorResponse{ServerErrorCode: CKErrorAuthenticationRequired, Reason: "auth required"}
	errBytes, _ := json.Marshal(errResp)

	mockedClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		return true
	})).Return(
		&http.Response{StatusCode: 401, Body: io.NopCloser(bytes.NewReader(errBytes))},
		nil,
	)

	result, err := ck.ListZones(CKDatabasePrivate)
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestLookupZones(t *testing.T) {
	mockedClient := new(MockedCKHTTPClient)
	ck := newTestCloudKit(mockedClient)

	resp := CKZonesResponse{
		Zones: []CKZone{
			{ZoneID: CKZoneID{ZoneName: "MyZone"}},
		},
	}
	respBytes, _ := json.Marshal(resp)

	mockedClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		return req.URL.Path == "/database/1/iCloud.com.example.app/development/private/zones/lookup"
	})).Return(
		&http.Response{StatusCode: 200, Body: io.NopCloser(bytes.NewReader(respBytes))},
		nil,
	)

	result, err := ck.LookupZones(CKDatabasePrivate, []CKZoneID{{ZoneName: "MyZone"}})
	assert.NoError(t, err)
	assert.Len(t, result.Zones, 1)
}

func TestModifyZones(t *testing.T) {
	mockedClient := new(MockedCKHTTPClient)
	ck := newTestCloudKit(mockedClient)

	resp := CKZonesResponse{
		Zones: []CKZone{
			{ZoneID: CKZoneID{ZoneName: "NewZone"}},
		},
	}
	respBytes, _ := json.Marshal(resp)

	mockedClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		return req.URL.Path == "/database/1/iCloud.com.example.app/development/private/zones/modify"
	})).Return(
		&http.Response{StatusCode: 200, Body: io.NopCloser(bytes.NewReader(respBytes))},
		nil,
	)

	result, err := ck.ModifyZones(CKDatabasePrivate, []CKZone{
		{ZoneID: CKZoneID{ZoneName: "NewZone"}},
	}, CKOperationCreate)
	assert.NoError(t, err)
	assert.Len(t, result.Zones, 1)
	assert.Equal(t, "NewZone", result.Zones[0].ZoneID.ZoneName)
}

func TestZoneChanges(t *testing.T) {
	mockedClient := new(MockedCKHTTPClient)
	ck := newTestCloudKit(mockedClient)

	resp := CKZoneChangesResponse{
		Zones: []CKZone{{ZoneID: CKZoneID{ZoneName: "ChangedZone"}}},
		SyncToken: "new-token",
	}
	respBytes, _ := json.Marshal(resp)

	mockedClient.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		return req.URL.Path == "/database/1/iCloud.com.example.app/development/private/zones/changes"
	})).Return(
		&http.Response{StatusCode: 200, Body: io.NopCloser(bytes.NewReader(respBytes))},
		nil,
	)

	result, err := ck.ZoneChanges(CKDatabasePrivate, &CKZoneChangesRequest{
		ZoneIDs: []CKZoneID{{ZoneName: "ChangedZone"}},
	})
	assert.NoError(t, err)
	assert.Equal(t, "new-token", result.SyncToken)
}
