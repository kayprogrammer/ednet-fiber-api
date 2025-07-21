package certs

import (
	"bytes"
	"fmt"
	"image/color"

	"github.com/fogleman/gg"
	"github.com/kayprogrammer/ednet-fiber-api/config"
	"github.com/kayprogrammer/ednet-fiber-api/ent"
)

const (
	width             = 1000
	height            = 700
	fontPath          = "DejaVuSans-Bold.ttf"
	signatureFontPath = "DejaVuSans.ttf"
)

func GenerateCertificate(student *ent.User, course *ent.Course, instructorName string) string {
	developerName := "Kenechi Ifeanyi"

	dc := gg.NewContext(width, height)

	dc.SetRGB(0.98, 0.98, 0.98)
	dc.Clear()

	dc.SetLineWidth(8)
	dc.SetRGB(0.2, 0.2, 0.6)
	dc.DrawRectangle(20, 20, float64(width-40), float64(height-40))
	dc.Stroke()

	_ = dc.LoadFontFace(fontPath, 48)
	dc.SetRGB(0.2, 0.3, 0.8)
	dc.DrawStringAnchored("ED", 80, 70, 0.5, 0.5)

	_ = dc.LoadFontFace(fontPath, 16)
	dc.SetRGB(0.1, 0.1, 0.1)
	dc.DrawStringAnchored("EDNET", 80, 110, 0.5, 0.5)

	_ = dc.LoadFontFace(fontPath, 48)
	dc.SetRGB(0.1, 0.1, 0.1)
	dc.DrawStringAnchored("~Certificate of Completion~", float64(width)/2, 120, 0.5, 0.5)

	_ = dc.LoadFontFace(fontPath, 36)
	dc.DrawStringAnchored("This is to certify that", float64(width)/2, 200, 0.5, 0.5)

	dc.SetRGB(0.1, 0.3, 0.7)
	dc.DrawStringAnchored(student.Name, float64(width)/2, 260, 0.5, 0.5)

	dc.SetRGB(0.1, 0.1, 0.1)
	dc.DrawStringAnchored("has successfully completed the course:", float64(width)/2, 320, 0.5, 0.5)

	dc.SetRGB(0.1, 0.3, 0.7)
	dc.DrawStringAnchored(course.Title, float64(width)/2, 380, 0.5, 0.5)

	dc.SetRGB(0.1, 0.1, 0.1)
	dc.DrawStringAnchored("at EDNET", float64(width)/2, 430, 0.5, 0.5)

	_ = dc.LoadFontFace(fontPath, 24)
	dc.SetRGB(0.1, 0.1, 0.1)
	dc.DrawStringAnchored("Developed By", 200, height-130, 0.5, 0.5)

	dc.SetLineWidth(3)
	dc.SetRGB(0.1, 0.1, 0.3)
	dc.DrawLine(100, height-111, 300, height-111)
	dc.Stroke()

	dc.SetRGB(0.2, 0.2, 0.2)
	_ = dc.LoadFontFace(fontPath, 22)
	dc.DrawStringAnchored(developerName, 200, height-90, 0.5, 0.5)

	_ = dc.LoadFontFace(fontPath, 24)
	dc.SetRGB(0.1, 0.1, 0.1)
	dc.DrawStringAnchored("Signature", 800, height-130, 0.5, 0.5)

	dc.SetLineWidth(2.5)
	dc.SetColor(color.Black)
	dc.DrawLine(700, height-111, 900, height-111)
	dc.Stroke()

	dc.SetRGB(0.2, 0.2, 0.2)
	_ = dc.LoadFontFace(fontPath, 22)
	dc.DrawStringAnchored(instructorName, 800, height-90, 0.5, 0.5)

	// Write to buffer instead of file
	buf := new(bytes.Buffer)
	if err := dc.EncodePNG(buf); err != nil {
		panic(err)
	}
	url := config.UploadGeneratedCert(buf, fmt.Sprintf("%s-certificate-for-%s", student.Username, course.Slug))
	return url
}