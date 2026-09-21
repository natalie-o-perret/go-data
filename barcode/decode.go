package barcode

import (
	"bytes"
	"errors"
	"image"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/aztec"
	"github.com/makiuchi-d/gozxing/datamatrix"
	"github.com/makiuchi-d/gozxing/oned"
	"github.com/makiuchi-d/gozxing/oned/rss"
	"github.com/makiuchi-d/gozxing/qrcode"
)

// ErrNotFound reports that an image contains no supported barcode.
var ErrNotFound = errors.New("barcode: not found")

// Result is decoded barcode data.
type Result struct {
	Text   string
	Bytes  []byte
	Format Format
}

// Decode finds and decodes the first supported barcode in img.
func Decode(img image.Image) (Result, error) {
	if img == nil {
		return Result{}, ErrNotFound
	}
	bitmap, err := gozxing.NewBinaryBitmapFromImage(img)
	if err != nil {
		return Result{}, err
	}
	hints := map[gozxing.DecodeHintType]interface{}{
		gozxing.DecodeHintType_TRY_HARDER:               true,
		gozxing.DecodeHintType_RETURN_CODABAR_START_END: true,
	}
	if result, ok := decodeWithHints(bitmap, hints); ok {
		return result, nil
	}
	hints[gozxing.DecodeHintType_PURE_BARCODE] = true
	if result, ok := decodeWithHints(bitmap, hints); ok {
		return result, nil
	}
	return Result{}, ErrNotFound
}

func decodeWithHints(bitmap *gozxing.BinaryBitmap,
	hints map[gozxing.DecodeHintType]interface{}) (Result, bool) {
	for _, reader := range newReaders() {
		result, err := reader.Decode(bitmap, hints)
		if err == nil {
			return Result{
				Text:   result.GetText(),
				Bytes:  bytes.Clone(result.GetRawBytes()),
				Format: Format(result.GetBarcodeFormat()),
			}, true
		}
	}
	return Result{}, false
}

func newReaders() []gozxing.Reader {
	return []gozxing.Reader{
		qrcode.NewQRCodeReader(),
		datamatrix.NewDataMatrixReader(),
		aztec.NewAztecReader(),
		oned.NewCode128Reader(),
		oned.NewCode93Reader(),
		oned.NewCode39Reader(),
		oned.NewCodaBarReader(),
		oned.NewUPCAReader(),
		oned.NewEAN8Reader(),
		oned.NewEAN13Reader(),
		oned.NewITFReader(),
		rss.NewRSS14Reader(),
	}
}
