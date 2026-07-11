package proxy

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image/png"

	"github.com/fogleman/gg"
	"golang.org/x/image/font/basicfont"
)

// CompressTextToImage renders a massive text payload into dense 1-bit PNG images (pages).
// It paginates the text to ensure no single image exceeds 2048px in height.
// Returns an array of base64 encoded PNG strings.
func CompressTextToImage(text string) ([]string, error) {
	var pages []string
	
	// Standard width for text rendering
	width := 1568
	charsPerLine := width / 7 // basicfont is 7x13
	
	// Max lines per page to keep height under ~2000px (2000 / 14 = ~142 lines)
	maxLinesPerPage := 140
	maxCharsPerPage := maxLinesPerPage * charsPerLine
	
	textLen := len(text)
	
	for i := 0; i < textLen; i += maxCharsPerPage {
		end := i + maxCharsPerPage
		if end > textLen {
			end = textLen
		}
		
		pageText := text[i:end]
		
		lines := (len(pageText) / charsPerLine) + 1
		height := (lines * 14) + 40

		dc := gg.NewContext(width, height)
		
		// Set background white
		dc.SetRGB(1, 1, 1)
		dc.Clear()
		
		// Set text color black
		dc.SetRGB(0, 0, 0)
		dc.SetFontFace(basicfont.Face7x13)
		
		dc.DrawStringWrapped(pageText, 10, 10, 0, 0, float64(width-20), 1.1, gg.AlignLeft)
		
		var buf bytes.Buffer
		if err := png.Encode(&buf, dc.Image()); err != nil {
			return nil, fmt.Errorf("failed to encode PNG: %w", err)
		}
		
		pages = append(pages, base64.StdEncoding.EncodeToString(buf.Bytes()))
	}
	
	return pages, nil
}
