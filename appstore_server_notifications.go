package apple

import "fmt"

// RequestTestNotification requests a test notification from the App Store Server.
func (s *appStoreServer) RequestTestNotification() (*ASTestNotificationResponse, error) {
	var result ASTestNotificationResponse
	if err := s.doRequest("POST", "/inApps/v1/notifications/test", nil, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetTestNotificationStatus gets the status of a test notification.
func (s *appStoreServer) GetTestNotificationStatus(testNotificationToken string) (*ASTestNotificationStatusResponse, error) {
	path := fmt.Sprintf("/inApps/v1/notifications/test/%s", testNotificationToken)
	var result ASTestNotificationStatusResponse
	if err := s.doRequest("GET", path, nil, nil, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetNotificationHistory gets the notification history.
func (s *appStoreServer) GetNotificationHistory(req *ASNotificationHistoryRequest) (*ASNotificationHistoryResponse, error) {
	var result ASNotificationHistoryResponse
	if err := s.doRequest("POST", "/inApps/v1/notifications/history", nil, req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
