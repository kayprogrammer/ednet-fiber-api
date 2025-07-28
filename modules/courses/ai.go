package courses

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/generative-ai-go/genai"
	"github.com/kayprogrammer/ednet-fiber-api/config"
	"github.com/ledongthuc/pdf"
	"google.golang.org/api/option"
)

func SummarizePDF(c *fiber.Ctx, cfg config.Config, maxPoints int) ([]string, int, *config.ErrorResponse) {
	apiKey := cfg.GeminiApiKey

	// Parse the multipart form:
	form, err := c.MultipartForm()
	errData := config.ErrorResponse{}
	if err != nil {
		errData = config.RequestErr(config.ERR_INVALID_ENTRY, "Error parsing form data")
		return nil, 422, &errData
	}

	// Get the file from the form
	files := form.File["file"]
	if len(files) == 0 {
		errData = config.RequestErr(config.ERR_INVALID_ENTRY, "File is required")
		log.Println("File is required")
		return nil, 422, &errData
	}
	file := files[0]

	// Check file size
	if file.Size > 3*1024*1024 {
		log.Println("File size cannot be more than 3MB")
		errData = config.RequestErr(config.ERR_INVALID_ENTRY, "File size cannot be more than 3MB")
		return nil, 422, &errData
	}

	// Read the PDF content
	textContent, err := readPDFContent(file)
	if err != nil {
		log.Printf("Error reading PDF content: %v", err)
		errData = config.RequestErr(config.ERR_INVALID_ENTRY, "Error reading PDF file")
		return nil, 422, &errData
	}

	ctx := c.Context()
	// Initialize the Gemini client
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Printf("Error initializing Gemini client: %v", err)
		errData = config.ServerErr("Error initializing Gemini client")
		return nil, 500, &errData
	}
	defer client.Close()

	// Choose the model
	model := client.GenerativeModel("gemini-2.0-flash")

	// Generate content
	prompt := []genai.Part{genai.Text(fmt.Sprintf(
		`
			Provide a detailed summary of the following text, 
			focusing strictly on the educational content and key concepts. 
			The summary should be a direct extraction of knowledge from the book, 
			presented as a list of bullet points. 
			Each point should comprehensively explain a specific concept, chapter, or key takeaway from the text. 
			Do not include any introductory phrases, metadata about the book (e.g., author, title, availability), or any other non-educational information. 
			The summary should have a maximum of %d points. Here is the text: %s
		`, maxPoints, textContent,
	))}
	resp, err := model.GenerateContent(ctx, prompt...)
	if err != nil {
		log.Printf("Error generating content from Gemini: %v", err)
		errData = config.ServerErr("Error generating content from Gemini")
		return nil, 500, &errData
	}

	// Format the response
	var summary string
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				summary += fmt.Sprintf("%v", part)
			}
		}
	}
	// Split the summary into a list of bullet points
	summaryPoints := strings.Split(summary, "\n")

	// Filter out any empty strings from the list
	var result []string
	for _, point := range summaryPoints {
		// Remove any asterisks and then trim the space
		cleanedPoint := strings.TrimSpace(strings.ReplaceAll(point, "*", ""))
		if cleanedPoint != "" {
			result = append(result, cleanedPoint)
		}
	}

	return result, 200, nil
}

func readPDFContent(fileHeader *multipart.FileHeader) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Create a temporary file to save the uploaded PDF
	tempFile, err := os.CreateTemp("", "uploaded-*.pdf")
	if err != nil {
		return "", err
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// Copy the uploaded file to the temporary file
	if _, err := io.Copy(tempFile, file); err != nil {
		return "", err
	}

	// Open the temporary file for reading
	pdfFile, err := os.Open(tempFile.Name())
	if err != nil {
		return "", err
	}
	defer pdfFile.Close()

	// Get file info for reader
	fileInfo, err := pdfFile.Stat()
	if err != nil {
		return "", err
	}

	// Create a new PDF reader
	pdfReader, err := pdf.NewReader(pdfFile, fileInfo.Size())
	if err != nil {
		return "", err
	}

	// Read all pages
	var textContent string
	numPages := pdfReader.NumPage()
	for i := 1; i <= numPages; i++ {
		page := pdfReader.Page(i)
		if page.V.IsNull() {
			continue
		}
		text, err := page.GetPlainText(nil)
		if err != nil {
			return "", err
		}
		textContent += text
	}

	return textContent, nil
}