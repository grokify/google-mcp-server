package slides

import (
	"context"
	"net/http"

	"google.golang.org/api/option"
	slides "google.golang.org/api/slides/v1"
)

// Service wraps the Google Slides API PresentationsService.
type Service struct {
	presentations *slides.PresentationsService
}

// NewService creates a new Slides service from an authenticated HTTP client.
func NewService(ctx context.Context, httpClient *http.Client) (*Service, error) {
	svc, err := slides.NewService(ctx, option.WithHTTPClient(httpClient))
	if err != nil {
		return nil, err
	}
	return &Service{
		presentations: slides.NewPresentationsService(svc),
	}, nil
}

// GetPresentation retrieves a presentation by ID.
func (s *Service) GetPresentation(ctx context.Context, presentationID string) (*slides.Presentation, error) {
	return s.presentations.Get(presentationID).Context(ctx).Do()
}
