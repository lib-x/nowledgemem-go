package nowledgemem

import "context"

// AdminService handles admin operations.
type AdminService struct {
	client *Client
}

// CheckUpgrade checks for available upgrades.
func (s *AdminService) CheckUpgrade(ctx context.Context) (*UpgradeInfo, error) {
	var resp UpgradeInfo
	if err := s.client.do(ctx, "GET", "/admin/upgrade/check", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// DownloadUpgrade downloads an available upgrade.
func (s *AdminService) DownloadUpgrade(ctx context.Context) error {
	return s.client.do(ctx, "POST", "/admin/upgrade/download", nil, nil)
}

// InstallUpgrade installs a downloaded upgrade.
func (s *AdminService) InstallUpgrade(ctx context.Context) error {
	return s.client.do(ctx, "POST", "/admin/upgrade/install", nil, nil)
}

// --- Admin types ---

// UpgradeInfo is the response for GET /admin/upgrade/check.
type UpgradeInfo struct {
	Available       bool   `json:"available"`
	CurrentVersion  string `json:"current_version"`
	LatestVersion   string `json:"latest_version,omitempty"`
	DownloadURL     string `json:"download_url,omitempty"`
	ReleaseNotes    string `json:"release_notes,omitempty"`
}
