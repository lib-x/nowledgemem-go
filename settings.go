package onledgemem

import "context"

// SettingsService handles settings operations.
type SettingsService struct {
	client *Client
}

// GetProfile returns user profile, aliases, context, and preferred language.
func (s *SettingsService) GetProfile(ctx context.Context) (*UserProfile, error) {
	var resp UserProfile
	if err := s.client.do(ctx, "GET", "/settings/profile", nil, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// --- Settings types ---

// UserProfile is the response for GET /settings/profile.
type UserProfile struct {
	Name              string `json:"name"`
	Aliases           string `json:"aliases"`
	Context           string `json:"context"`
	PreferredLanguage string `json:"preferred_language"`
	CustomInstructions string `json:"custom_instructions"`
}
