package barcode_test

import (
	"errors"
	"image"
	"testing"

	"github.com/natalie-o-perret/go-data/barcode"
)

func TestRoundTrip(t *testing.T) {
	tests := []struct {
		format        barcode.Format
		input, want   string
		width, height int
	}{
		{barcode.Codabar, "A1234B", "A1234B", 240, 80},
		{barcode.Code39, "HELLO-39", "HELLO-39", 240, 80},
		{barcode.Code93, "HELLO93", "HELLO93", 240, 80},
		{barcode.Code128, "HELLO-1D", "HELLO-1D", 240, 80},
		{barcode.DataMatrix, "hello matrix", "hello matrix", 160, 160},
		{barcode.EAN8, "1234567", "12345670", 240, 80},
		{barcode.EAN13, "123456789012", "1234567890128", 240, 80},
		{barcode.ITF, "123456", "123456", 240, 80},
		{barcode.QRCode, "hello 2d", "hello 2d", 160, 160},
		{barcode.UPCA, "12345678901", "123456789012", 240, 80},
	}
	for _, test := range tests {
		t.Run(test.format.String(), func(t *testing.T) {
			img, err := barcode.Encode(test.input, test.format, test.width, test.height)
			if err != nil {
				t.Fatal(err)
			}
			got, err := barcode.Decode(img)
			if err != nil {
				t.Fatal(err)
			}
			if got.Text != test.want || got.Format != test.format {
				t.Fatalf("Decode() = %#v, want text %q in %s", got, test.want, test.format)
			}
		})
	}
}

func TestFormatDimensions(t *testing.T) {
	if Code128, QRCode := barcode.Code128.Dimension(), barcode.QRCode.Dimension(); Code128 != 1 || QRCode != 2 {
		t.Fatalf("dimensions = (%d, %d), want (1, 2)", Code128, QRCode)
	}
}

func TestErrors(t *testing.T) {
	if _, err := barcode.Encode("x", barcode.Aztec, 100, 100); !errors.Is(err, barcode.ErrUnsupportedFormat) {
		t.Fatalf("Encode(Aztec) error = %v, want ErrUnsupportedFormat", err)
	}
	if _, err := barcode.Decode(image.NewGray(image.Rect(0, 0, 32, 32))); !errors.Is(err, barcode.ErrNotFound) {
		t.Fatalf("Decode(blank) error = %v, want ErrNotFound", err)
	}
}
