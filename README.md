# Google MCP Server

An MCP (Model Context Protocol) server for reading Google Slides presentations.

## Features

Read-only access to Google Slides presentations via five MCP tools:

- **get_presentation** - Get presentation metadata (title, slide count, locale, revision ID)
- **list_slides** - List all slides with titles and element counts
- **get_slide** - Get slide content and element details by index or object ID
- **get_slide_notes** - Get speaker notes by slide index or object ID
- **get_presentation_content** - Get all slides' text and images in one call (ideal for AI)

## Requirements

- Go 1.24+
- Google Cloud service account with Slides API access

## Installation

```bash
go install github.com/grokify/google-mcp-server/cmd/google-mcp-server@latest
```

Or build from source:

```bash
git clone https://github.com/grokify/google-mcp-server.git
cd google-mcp-server
go build ./cmd/google-mcp-server
```

## Setup

### 1. Create a Google Cloud Service Account

1. Go to the [Google Cloud Console](https://console.cloud.google.com/)
2. Create a new project or select an existing one
3. Enable the Google Slides API
4. Create a service account with no special roles
5. Download the JSON credentials file

### 2. Share Presentations with the Service Account

Share any presentations you want to access with the service account's email address (found in the credentials JSON as `client_email`).

## Usage

### Option 1: Google Service Account Credentials

Use a standard Google Cloud service account JSON file:

```bash
google-mcp-server -credentials /path/to/service-account.json
```

Or using an environment variable:

```bash
export GOOGLE_CREDENTIALS_FILE=/path/to/service-account.json
google-mcp-server
```

### Option 2: goauth CredentialsSet

Use a [goauth](https://github.com/grokify/goauth) CredentialsSet file, which can store multiple credentials:

```bash
google-mcp-server -goauth-credentials-file /path/to/credentials.json -goauth-credentials-account myaccount
```

Or using environment variables:

```bash
export GOAUTH_CREDENTIALS_FILE=/path/to/credentials.json
export GOAUTH_CREDENTIALS_ACCOUNT=myaccount
google-mcp-server
```

The CredentialsSet entry should be of type `gcpsa` with appropriate scopes:

```json
{
  "credentials": {
    "myaccount": {
      "type": "gcpsa",
      "gcpsa": {
        "gcpCredentials": {
          "type": "service_account",
          "project_id": "...",
          "private_key_id": "...",
          "private_key": "...",
          "client_email": "...",
          "client_id": "..."
        },
        "scopes": [
          "https://www.googleapis.com/auth/presentations.readonly",
          "https://www.googleapis.com/auth/drive.readonly"
        ]
      }
    }
  }
}
```

### Claude Desktop Configuration

Add to your Claude Desktop configuration (`claude_desktop_config.json`):

```json
{
  "mcpServers": {
    "google": {
      "command": "/path/to/google-mcp-server",
      "args": ["-credentials", "/path/to/service-account.json"]
    }
  }
}
```

Or with goauth:

```json
{
  "mcpServers": {
    "google": {
      "command": "/path/to/google-mcp-server",
      "args": [
        "-goauth-credentials-file", "/path/to/credentials.json",
        "-goauth-credentials-account", "myaccount"
      ]
    }
  }
}
```

## Tools

### get_presentation

Get metadata about a presentation.

**Input:**

- `presentation_id` (required) - The ID of the Google Slides presentation

**Output:**

- `title` - Presentation title
- `slide_count` - Number of slides
- `locale` - Presentation locale
- `revision_id` - Current revision ID

### list_slides

List all slides in a presentation.

**Input:**

- `presentation_id` (required) - The ID of the Google Slides presentation

**Output:**

- `slides` - Array of slide information:
  - `object_id` - Slide's unique identifier
  - `index` - Zero-based slide index
  - `title` - Slide title (if present)
  - `element_count` - Number of elements on the slide

### get_slide

Get the content of a specific slide.

**Input:**

- `presentation_id` (required) - The ID of the Google Slides presentation
- `slide_index` (optional) - Zero-based slide index
- `slide_object_id` (optional) - Slide's object ID

One of `slide_index` or `slide_object_id` must be provided.

**Output:**

- `text_content` - Array of text strings from the slide
- `element_summary` - Array of element details:
  - `object_id` - Element's unique identifier
  - `element_type` - Type of element (shape, image, table, etc.)
  - `description` - Element description or text preview

### get_slide_notes

Get the speaker notes for a specific slide.

**Input:**

- `presentation_id` (required) - The ID of the Google Slides presentation
- `slide_index` (optional) - Zero-based slide index
- `slide_object_id` (optional) - Slide's object ID

One of `slide_index` or `slide_object_id` must be provided.

**Output:**

- `notes` - Speaker notes text

### get_presentation_content

Get all slide content in a single call - ideal for AI analysis of the entire presentation.

**Input:**

- `presentation_id` (required) - The ID of the Google Slides presentation
- `include_notes` (optional) - Include speaker notes for each slide (default: false)

**Output:**

- `title` - Presentation title
- `slides` - Array of slide content:
  - `index` - Zero-based slide index
  - `object_id` - Slide's unique identifier
  - `title` - Slide title (if present)
  - `text_content` - Array of text strings from the slide
  - `images` - Array of images:
    - `object_id` - Image element ID
    - `content_url` - Direct URL to image (valid ~30 minutes)
    - `source_url` - Original source URL (if available)
    - `alt_text` - Image description
  - `notes` - Speaker notes (if `include_notes` is true)

## Finding Presentation IDs

The presentation ID is the long string in the URL when viewing a presentation:

```
https://docs.google.com/presentation/d/PRESENTATION_ID_HERE/edit
```

## License

MIT
