// Package slides provides MCP tools for reading Google Slides presentations.
package slides

// GetPresentationInput is the input for the get_presentation tool.
type GetPresentationInput struct {
	PresentationID string `json:"presentation_id" jsonschema:"required,description=The ID of the Google Slides presentation"`
}

// GetPresentationOutput is the output for the get_presentation tool.
type GetPresentationOutput struct {
	Title      string `json:"title"`
	SlideCount int    `json:"slide_count"`
	Locale     string `json:"locale,omitempty"`
	RevisionID string `json:"revision_id,omitempty"`
}

// ListSlidesInput is the input for the list_slides tool.
type ListSlidesInput struct {
	PresentationID string `json:"presentation_id" jsonschema:"required,description=The ID of the Google Slides presentation"`
}

// SlideInfo represents basic information about a slide.
type SlideInfo struct {
	ObjectID     string `json:"object_id"`
	Index        int    `json:"index"`
	Title        string `json:"title,omitempty"`
	ElementCount int    `json:"element_count"`
}

// ListSlidesOutput is the output for the list_slides tool.
type ListSlidesOutput struct {
	Slides []SlideInfo `json:"slides"`
}

// GetSlideInput is the input for the get_slide tool.
type GetSlideInput struct {
	PresentationID string `json:"presentation_id" jsonschema:"required,description=The ID of the Google Slides presentation"`
	SlideIndex     *int   `json:"slide_index,omitempty" jsonschema:"description=The zero-based index of the slide (mutually exclusive with slide_object_id)"`
	SlideObjectID  string `json:"slide_object_id,omitempty" jsonschema:"description=The object ID of the slide (mutually exclusive with slide_index)"`
}

// ElementSummary represents a summary of a page element.
type ElementSummary struct {
	ObjectID    string `json:"object_id"`
	ElementType string `json:"element_type"`
	Description string `json:"description,omitempty"`
	ImageURL    string `json:"image_url,omitempty"`
}

// GetSlideOutput is the output for the get_slide tool.
type GetSlideOutput struct {
	TextContent    []string         `json:"text_content"`
	ElementSummary []ElementSummary `json:"element_summary"`
	Images         []ImageInfo      `json:"images,omitempty"`
}

// ImageInfo represents an image on a slide.
type ImageInfo struct {
	ObjectID   string `json:"object_id"`
	ContentURL string `json:"content_url"`
	SourceURL  string `json:"source_url,omitempty"`
	AltText    string `json:"alt_text,omitempty"`
}

// GetSlideNotesInput is the input for the get_slide_notes tool.
type GetSlideNotesInput struct {
	PresentationID string `json:"presentation_id" jsonschema:"required,description=The ID of the Google Slides presentation"`
	SlideIndex     *int   `json:"slide_index,omitempty" jsonschema:"description=The zero-based index of the slide (mutually exclusive with slide_object_id)"`
	SlideObjectID  string `json:"slide_object_id,omitempty" jsonschema:"description=The object ID of the slide (mutually exclusive with slide_index)"`
}

// GetSlideNotesOutput is the output for the get_slide_notes tool.
type GetSlideNotesOutput struct {
	Notes string `json:"notes"`
}

// GetPresentationContentInput is the input for the get_presentation_content tool.
type GetPresentationContentInput struct {
	PresentationID string `json:"presentation_id" jsonschema:"required,description=The ID of the Google Slides presentation"`
	IncludeNotes   bool   `json:"include_notes,omitempty" jsonschema:"description=Include speaker notes for each slide"`
}

// SlideContent represents the content of a single slide.
type SlideContent struct {
	Index       int         `json:"index"`
	ObjectID    string      `json:"object_id"`
	Title       string      `json:"title,omitempty"`
	TextContent []string    `json:"text_content"`
	Images      []ImageInfo `json:"images,omitempty"`
	Notes       string      `json:"notes,omitempty"`
}

// GetPresentationContentOutput is the output for the get_presentation_content tool.
type GetPresentationContentOutput struct {
	Title  string         `json:"title"`
	Slides []SlideContent `json:"slides"`
}
