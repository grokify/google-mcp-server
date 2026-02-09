package slides

import (
	"context"
	"fmt"

	"github.com/grokify/gogoogle/slidesutil/v1"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Tools wraps the Slides service and provides MCP tool handlers.
type Tools struct {
	service *Service
}

// NewTools creates a new Tools instance.
func NewTools(service *Service) *Tools {
	return &Tools{service: service}
}

// GetPresentation retrieves metadata about a presentation.
func (t *Tools) GetPresentation(ctx context.Context, req *mcp.CallToolRequest, input GetPresentationInput) (*mcp.CallToolResult, GetPresentationOutput, error) {
	pres, err := t.service.GetPresentation(ctx, input.PresentationID)
	if err != nil {
		return nil, GetPresentationOutput{}, fmt.Errorf("failed to get presentation: %w", err)
	}

	output := GetPresentationOutput{
		Title:      pres.Title,
		SlideCount: len(pres.Slides),
		Locale:     pres.Locale,
		RevisionID: pres.RevisionId,
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: fmt.Sprintf("Presentation: %s\nSlides: %d\nLocale: %s\nRevision: %s",
				output.Title, output.SlideCount, output.Locale, output.RevisionID)},
		},
	}, output, nil
}

// ListSlides returns a list of all slides in a presentation.
func (t *Tools) ListSlides(ctx context.Context, req *mcp.CallToolRequest, input ListSlidesInput) (*mcp.CallToolResult, ListSlidesOutput, error) {
	pres, err := t.service.GetPresentation(ctx, input.PresentationID)
	if err != nil {
		return nil, ListSlidesOutput{}, fmt.Errorf("failed to get presentation: %w", err)
	}

	slides := make([]SlideInfo, len(pres.Slides))
	for i, slide := range pres.Slides {
		slides[i] = SlideInfo{
			ObjectID:     slide.ObjectId,
			Index:        i,
			Title:        extractSlideTitle(slide),
			ElementCount: len(slide.PageElements),
		}
	}

	output := ListSlidesOutput{Slides: slides}

	// Build text summary
	var summary string
	for _, s := range slides {
		title := s.Title
		if title == "" {
			title = "(untitled)"
		}
		summary += fmt.Sprintf("Slide %d [%s]: %s (%d elements)\n", s.Index, s.ObjectID, title, s.ElementCount)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: summary},
		},
	}, output, nil
}

// GetSlide retrieves the content of a specific slide.
func (t *Tools) GetSlide(ctx context.Context, req *mcp.CallToolRequest, input GetSlideInput) (*mcp.CallToolResult, GetSlideOutput, error) {
	pres, err := t.service.GetPresentation(ctx, input.PresentationID)
	if err != nil {
		return nil, GetSlideOutput{}, fmt.Errorf("failed to get presentation: %w", err)
	}

	slide, idx, err := findSlide(pres, input.SlideIndex, input.SlideObjectID)
	if err != nil {
		return nil, GetSlideOutput{}, err
	}

	textContent := extractTextContent(slide)
	images := extractImages(slide)
	elements := make([]ElementSummary, len(slide.PageElements))
	for i, elem := range slide.PageElements {
		elements[i] = ElementSummary{
			ObjectID:    elem.ObjectId,
			ElementType: getElementType(elem),
			Description: getElementDescription(elem),
			ImageURL:    getImageURL(elem),
		}
	}

	output := GetSlideOutput{
		TextContent:    textContent,
		ElementSummary: elements,
		Images:         images,
	}

	// Build text summary
	title := extractSlideTitle(slide)
	if title == "" {
		title = "(untitled)"
	}
	summary := fmt.Sprintf("Slide %d: %s\n\n", idx, title)
	summary += "Text content:\n"
	for _, text := range textContent {
		summary += fmt.Sprintf("  - %s\n", text)
	}
	if len(images) > 0 {
		summary += fmt.Sprintf("\nImages (%d):\n", len(images))
		for _, img := range images {
			alt := img.AltText
			if alt == "" {
				alt = "(no alt text)"
			}
			summary += fmt.Sprintf("  - [%s] %s\n    URL: %s\n", img.ObjectID, alt, img.ContentURL)
		}
	}
	summary += "\nElements:\n"
	for _, elem := range elements {
		desc := elem.Description
		if desc != "" {
			desc = ": " + desc
		}
		summary += fmt.Sprintf("  - [%s] %s%s\n", elem.ObjectID, elem.ElementType, desc)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: summary},
		},
	}, output, nil
}

// GetSlideNotes retrieves the speaker notes for a specific slide.
func (t *Tools) GetSlideNotes(ctx context.Context, req *mcp.CallToolRequest, input GetSlideNotesInput) (*mcp.CallToolResult, GetSlideNotesOutput, error) {
	pres, err := t.service.GetPresentation(ctx, input.PresentationID)
	if err != nil {
		return nil, GetSlideNotesOutput{}, fmt.Errorf("failed to get presentation: %w", err)
	}

	slide, idx, err := findSlide(pres, input.SlideIndex, input.SlideObjectID)
	if err != nil {
		return nil, GetSlideNotesOutput{}, err
	}

	notes := extractNotesText(slide)
	output := GetSlideNotesOutput{Notes: notes}

	summary := fmt.Sprintf("Speaker notes for slide %d:\n\n%s", idx, notes)
	if notes == "" {
		summary = fmt.Sprintf("Slide %d has no speaker notes.", idx)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: summary},
		},
	}, output, nil
}

// GetPresentationContent retrieves all slide content in a single call.
func (t *Tools) GetPresentationContent(ctx context.Context, req *mcp.CallToolRequest, input GetPresentationContentInput) (*mcp.CallToolResult, GetPresentationContentOutput, error) {
	pres, err := t.service.GetPresentation(ctx, input.PresentationID)
	if err != nil {
		return nil, GetPresentationContentOutput{}, fmt.Errorf("failed to get presentation: %w", err)
	}

	// Use gogoogle's ExtractPresentationContent
	content := slidesutil.ExtractPresentationContent(pres, input.IncludeNotes)

	// Convert to MCP output types
	slides := make([]SlideContent, len(content.Slides))
	for i, sc := range content.Slides {
		images := make([]ImageInfo, len(sc.Images))
		for j, img := range sc.Images {
			images[j] = ImageInfo{
				ObjectID:   img.ObjectID,
				ContentURL: img.ContentURL,
				SourceURL:  img.SourceURL,
				AltText:    img.AltText,
			}
		}
		slides[i] = SlideContent{
			Index:       sc.Index,
			ObjectID:    sc.ObjectID,
			Title:       sc.Title,
			TextContent: sc.TextContent,
			Images:      images,
			Notes:       sc.Notes,
		}
	}

	output := GetPresentationContentOutput{
		Title:  content.Title,
		Slides: slides,
	}

	// Build text summary
	summary := fmt.Sprintf("Presentation: %s (%d slides)\n\n", content.Title, len(slides))
	for _, slide := range slides {
		title := slide.Title
		if title == "" {
			title = "(untitled)"
		}
		summary += fmt.Sprintf("--- Slide %d: %s ---\n", slide.Index, title)
		if len(slide.TextContent) > 0 {
			for _, text := range slide.TextContent {
				summary += fmt.Sprintf("  %s\n", text)
			}
		}
		if len(slide.Images) > 0 {
			summary += fmt.Sprintf("  [%d image(s)]\n", len(slide.Images))
		}
		if slide.Notes != "" {
			summary += fmt.Sprintf("  Notes: %s\n", slide.Notes)
		}
		summary += "\n"
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: summary},
		},
	}, output, nil
}

// RegisterTools registers all Slides tools with the MCP server.
func RegisterTools(server *mcp.Server, tools *Tools) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_presentation",
		Description: "Get metadata about a Google Slides presentation including title, slide count, locale, and revision ID",
	}, tools.GetPresentation)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_slides",
		Description: "List all slides in a Google Slides presentation with their titles and element counts",
	}, tools.ListSlides)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_slide",
		Description: "Get the content and elements of a specific slide by index or object ID",
	}, tools.GetSlide)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_slide_notes",
		Description: "Get the speaker notes for a specific slide by index or object ID",
	}, tools.GetSlideNotes)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_presentation_content",
		Description: "Get all slide content (text and images) in a single call, ideal for AI analysis of the entire presentation",
	}, tools.GetPresentationContent)
}
