package apple

// CKEnvironment represents the CloudKit environment.
type CKEnvironment string

const (
	CKEnvironmentDevelopment CKEnvironment = "development"
	CKEnvironmentProduction  CKEnvironment = "production"
)

// CKDatabase represents the CloudKit database type.
type CKDatabase string

const (
	CKDatabasePublic  CKDatabase = "public"
	CKDatabasePrivate CKDatabase = "private"
	CKDatabaseShared  CKDatabase = "shared"
)

// CKOperationType represents the type of operation to perform on a record.
type CKOperationType string

const (
	CKOperationCreate       CKOperationType = "create"
	CKOperationUpdate       CKOperationType = "update"
	CKOperationForceUpdate  CKOperationType = "forceUpdate"
	CKOperationReplace      CKOperationType = "replace"
	CKOperationForceReplace CKOperationType = "forceReplace"
	CKOperationDelete       CKOperationType = "delete"
	CKOperationForceDelete  CKOperationType = "forceDelete"
)

// CKField represents a field in a CloudKit record.
type CKField struct {
	Value any `json:"value"`
	Type  string      `json:"type,omitempty"`
}

// CKRecord represents a CloudKit record.
type CKRecord struct {
	RecordName       string              `json:"recordName,omitempty"`
	RecordType       string              `json:"recordType,omitempty"`
	RecordChangeTag  string              `json:"recordChangeTag,omitempty"`
	Fields           map[string]*CKField `json:"fields,omitempty"`
	Created          *CKTimestamp        `json:"created,omitempty"`
	Modified         *CKTimestamp        `json:"modified,omitempty"`
	Deleted          bool                `json:"deleted,omitempty"`
	ZoneID           *CKZoneID           `json:"zoneID,omitempty"`
	PluginFields     map[string]*CKField `json:"pluginFields,omitempty"`
	DesiredKeys      []string            `json:"desiredKeys,omitempty"`
	NumbersAsStrings bool                `json:"numbersAsStrings,omitempty"`
}

// CKTimestamp represents a CloudKit timestamp with user and device info.
type CKTimestamp struct {
	Timestamp    int64  `json:"timestamp,omitempty"`
	UserRecordName string `json:"userRecordName,omitempty"`
	DeviceID     string `json:"deviceID,omitempty"`
}

// CKZoneID represents a CloudKit zone identifier.
type CKZoneID struct {
	ZoneName  string `json:"zoneName"`
	OwnerName string `json:"ownerRecordName,omitempty"`
}

// CKZone represents a CloudKit zone.
type CKZone struct {
	ZoneID    CKZoneID `json:"zoneID"`
	SyncToken string   `json:"syncToken,omitempty"`
	Atomic    bool     `json:"atomic,omitempty"`
}

// CKQuery represents a query for CloudKit records.
type CKQuery struct {
	RecordType  string      `json:"recordType"`
	FilterBy    []CKFilter  `json:"filterBy,omitempty"`
	SortBy      []CKSort    `json:"sortBy,omitempty"`
}

// CKFilter represents a filter in a CloudKit query.
type CKFilter struct {
	Comparator string      `json:"comparator"`
	FieldName  string      `json:"fieldName"`
	FieldValue *CKField    `json:"fieldValue"`
	Distance   float64     `json:"distance,omitempty"`
}

// CKSort represents a sort descriptor in a CloudKit query.
type CKSort struct {
	FieldName string `json:"fieldName"`
	Ascending bool   `json:"ascending"`
}

// CKRecordOperation represents an operation on a CloudKit record.
type CKRecordOperation struct {
	OperationType CKOperationType `json:"operationType"`
	Record        CKRecord        `json:"record"`
	DesiredKeys   []string        `json:"desiredKeys,omitempty"`
}

// CKSubscription represents a CloudKit subscription.
type CKSubscription struct {
	SubscriptionID   string            `json:"subscriptionID,omitempty"`
	SubscriptionType string            `json:"subscriptionType,omitempty"`
	Query            *CKQuery          `json:"query,omitempty"`
	FiresOn          []string          `json:"firesOn,omitempty"`
	FiresOnce        bool              `json:"firesOnce,omitempty"`
	NotificationInfo *CKNotificationInfo `json:"notificationInfo,omitempty"`
	ZoneID           *CKZoneID         `json:"zoneID,omitempty"`
	ZoneWide         bool              `json:"zoneWide,omitempty"`
}

// CKSubscriptionOperation represents an operation on a CloudKit subscription.
type CKSubscriptionOperation struct {
	OperationType CKOperationType `json:"operationType"`
	Subscription  CKSubscription  `json:"subscription"`
}

// CKNotificationInfo represents push notification configuration.
type CKNotificationInfo struct {
	AlertBody             string   `json:"alertBody,omitempty"`
	AlertLocalizationKey  string   `json:"alertLocalizationKey,omitempty"`
	AlertLocalizationArgs []string `json:"alertLocalizationArgs,omitempty"`
	AlertActionLocalizationKey string `json:"alertActionLocalizationKey,omitempty"`
	AlertLaunchImage      string   `json:"alertLaunchImage,omitempty"`
	SoundName             string   `json:"soundName,omitempty"`
	DesiredKeys           []string `json:"desiredKeys,omitempty"`
	ShouldBadge           bool     `json:"shouldBadge,omitempty"`
	ShouldSendContentAvailable bool `json:"shouldSendContentAvailable,omitempty"`
	ShouldSendMutableContent   bool `json:"shouldSendMutableContent,omitempty"`
}

// CKAsset represents a CloudKit asset.
type CKAsset struct {
	FileChecksum    string `json:"fileChecksum,omitempty"`
	Size            int64  `json:"size,omitempty"`
	DownloadURL     string `json:"downloadURL,omitempty"`
	ReferenceChecksum string `json:"referenceChecksum,omitempty"`
	WrappingKey     string `json:"wrappingKey,omitempty"`
	Receipt         string `json:"receipt,omitempty"`
}

// CKAssetUploadRequest represents a request to upload CloudKit assets.
type CKAssetUploadRequest struct {
	RecordType string `json:"recordType"`
	FieldName  string `json:"fieldName"`
}

// CKAssetUploadResponse represents a response for asset upload.
type CKAssetUploadResponse struct {
	RecordName string `json:"recordName,omitempty"`
	FieldName  string `json:"fieldName,omitempty"`
	URL        string `json:"url,omitempty"`
}

// CKUserInfo represents a CloudKit user identity.
type CKUserInfo struct {
	UserRecordName  string `json:"userRecordName,omitempty"`
	FirstName       string `json:"firstName,omitempty"`
	LastName        string `json:"lastName,omitempty"`
	EmailAddress    string `json:"emailAddress,omitempty"`
	IsDiscoverable  bool   `json:"isDiscoverable,omitempty"`
}

// CKAPNsToken represents an APNs token for CloudKit push notifications.
type CKAPNsToken struct {
	APNsToken      string `json:"apnsToken"`
	APNsEnvironment string `json:"apnsEnvironment,omitempty"`
}

// CKQueryRequest represents a request to query CloudKit records.
type CKQueryRequest struct {
	Query            CKQuery   `json:"query"`
	ZoneID           *CKZoneID `json:"zoneID,omitempty"`
	ResultsLimit     int       `json:"resultsLimit,omitempty"`
	DesiredKeys      []string  `json:"desiredKeys,omitempty"`
	ContinuationMarker string  `json:"continuationMarker,omitempty"`
	NumbersAsStrings bool      `json:"numbersAsStrings,omitempty"`
	ZoneWide         bool      `json:"zoneWide,omitempty"`
}

// CKQueryResponse represents a response from a CloudKit query.
type CKQueryResponse struct {
	Records            []CKRecord `json:"records"`
	ContinuationMarker string     `json:"continuationMarker,omitempty"`
}

// CKRecordsModifyRequest represents a request to modify CloudKit records.
type CKRecordsModifyRequest struct {
	Operations       []CKRecordOperation `json:"operations"`
	ZoneID           *CKZoneID           `json:"zoneID,omitempty"`
	Atomic           bool                `json:"atomic,omitempty"`
	DesiredKeys      []string            `json:"desiredKeys,omitempty"`
	NumbersAsStrings bool                `json:"numbersAsStrings,omitempty"`
}

// CKRecordsModifyResponse represents a response from modifying CloudKit records.
type CKRecordsModifyResponse struct {
	Records []CKRecord `json:"records"`
}

// CKRecordsLookupRequest represents a request to look up CloudKit records.
type CKRecordsLookupRequest struct {
	Records          []CKRecord `json:"records"`
	ZoneID           *CKZoneID  `json:"zoneID,omitempty"`
	DesiredKeys      []string   `json:"desiredKeys,omitempty"`
	NumbersAsStrings bool       `json:"numbersAsStrings,omitempty"`
}

// CKRecordsLookupResponse represents a response from looking up CloudKit records.
type CKRecordsLookupResponse struct {
	Records []CKRecord `json:"records"`
}

// CKRecordChangesRequest represents a request to get record changes.
type CKRecordChangesRequest struct {
	ZoneID           CKZoneID `json:"zoneID"`
	SyncToken        string   `json:"syncToken,omitempty"`
	ResultsLimit     int      `json:"resultsLimit,omitempty"`
	DesiredKeys      []string `json:"desiredKeys,omitempty"`
	NumbersAsStrings bool     `json:"numbersAsStrings,omitempty"`
}

// CKRecordChangesResponse represents a response from getting record changes.
type CKRecordChangesResponse struct {
	Records   []CKRecord `json:"records"`
	SyncToken string     `json:"syncToken,omitempty"`
	MoreComing bool      `json:"moreComing,omitempty"`
}

// CKZonesResponse represents a response containing CloudKit zones.
type CKZonesResponse struct {
	Zones []CKZone `json:"zones"`
}

// CKZoneChangesRequest represents a request to get zone changes.
type CKZoneChangesRequest struct {
	ZoneIDs []CKZoneID `json:"zones"`
	SyncToken string   `json:"syncToken,omitempty"`
}

// CKZoneChangesResponse represents a response from getting zone changes.
type CKZoneChangesResponse struct {
	Zones     []CKZone `json:"zones"`
	SyncToken string   `json:"syncToken,omitempty"`
	MoreComing bool    `json:"moreComing,omitempty"`
}

// CKSubscriptionsResponse represents a response containing CloudKit subscriptions.
type CKSubscriptionsResponse struct {
	Subscriptions []CKSubscription `json:"subscriptions"`
}

// CKSubscriptionsModifyRequest represents a request to modify CloudKit subscriptions.
type CKSubscriptionsModifyRequest struct {
	Operations []CKSubscriptionOperation `json:"operations"`
}

// CKSubscriptionsModifyResponse represents a response from modifying CloudKit subscriptions.
type CKSubscriptionsModifyResponse struct {
	Subscriptions []CKSubscription `json:"subscriptions"`
}

// CKAssetsUploadRequest represents a request to upload CloudKit assets.
type CKAssetsUploadRequest struct {
	Tokens []CKAssetUploadRequest `json:"tokens"`
	ZoneID *CKZoneID              `json:"zoneID,omitempty"`
}

// CKAssetsUploadResponse represents a response from uploading CloudKit assets.
type CKAssetsUploadResponse struct {
	Tokens []CKAssetUploadResponse `json:"tokens"`
}

// CKUsersResponse represents a response containing CloudKit users.
type CKUsersResponse struct {
	Users []CKUserInfo `json:"users"`
}

// CKUserLookupRequest represents a request to look up CloudKit users.
type CKUserLookupRequest struct {
	EmailAddresses []string `json:"emailAddresses,omitempty"`
	PhoneNumbers   []string `json:"phoneNumbers,omitempty"`
}

// CKTokensCreateRequest represents a request to create CloudKit APNs tokens.
type CKTokensCreateRequest struct {
	Tokens []CKAPNsToken `json:"tokens"`
}

// CKTokensCreateResponse represents a response from creating CloudKit APNs tokens.
type CKTokensCreateResponse struct {
	Tokens []CKAPNsToken `json:"tokens"`
}
