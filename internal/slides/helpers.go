package slides

import (
	"github.com/grokify/gogoogle/slidesutil/v1"
	gslides "google.golang.org/api/slides/v1"
)

// findSlide locates a slide by either index or object ID.
func findSlide(pres *gslides.Presentation, index *int, objectID string) (*gslides.Page, int, error) {
	return slidesutil.FindSlide(pres, index, objectID)
}

// extractSlideTitle extracts the title text from a slide.
func extractSlideTitle(slide *gslides.Page) string {
	return slidesutil.ExtractSlideTitle(slide)
}

// extractTextContent extracts all text content from a page's elements.
func extractTextContent(page *gslides.Page) []string {
	return slidesutil.ExtractTextContent(page)
}

// extractNotesText extracts speaker notes text from a slide.
func extractNotesText(slide *gslides.Page) string {
	return slidesutil.ExtractNotesText(slide)
}

// getElementType returns a human-readable type name for a page element.
func getElementType(elem *gslides.PageElement) string {
	return slidesutil.GetElementType(elem)
}

// getElementDescription returns a description for a page element.
func getElementDescription(elem *gslides.PageElement) string {
	return slidesutil.GetElementDescription(elem)
}

// getImageURL returns the content URL for an image element, if any.
func getImageURL(elem *gslides.PageElement) string {
	return slidesutil.GetImageURL(elem)
}

// extractImages extracts all images from a page and converts to MCP ImageInfo.
func extractImages(page *gslides.Page) []ImageInfo {
	gogoogleImages := slidesutil.ExtractImages(page)
	images := make([]ImageInfo, len(gogoogleImages))
	for i, img := range gogoogleImages {
		images[i] = ImageInfo{
			ObjectID:   img.ObjectID,
			ContentURL: img.ContentURL,
			SourceURL:  img.SourceURL,
			AltText:    img.AltText,
		}
	}
	return images
}
