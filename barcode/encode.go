package barcode

import (
	"errors"
	"fmt"
	"image"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/datamatrix"
	"github.com/makiuchi-d/gozxing/oned"
	"github.com/makiuchi-d/gozxing/qrcode"
)

// ErrUnsupportedFormat reports a format that cannot be encoded.
var ErrUnsupportedFormat = errors.New("barcode: unsupported format")

// Encode returns an image containing contents in the requested format.
func Encode(contents string, format Format, width, height int) (image.Image, error) {
	if width <= 0 || height <= 0 {
		return nil, errors.New("barcode: width and height must be positive")
	}
	writer := newWriter(format)
	if writer == nil {
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedFormat, format)
	}
	return writer.EncodeWithoutHint(contents, gozxing.BarcodeFormat(format), width, height)
}

func newWriter(format Format) gozxing.Writer {
	switch format {
	case Codabar:
		return oned.NewCodaBarWriter()
	case Code39:
		return oned.NewCode39Writer()
	case Code93:
		return oned.NewCode93Writer()
	case Code128:
		return oned.NewCode128Writer()
	case DataMatrix:
		return datamatrix.NewDataMatrixWriter()
	case EAN8:
		return oned.NewEAN8Writer()
	case EAN13:
		return oned.NewEAN13Writer()
	case ITF:
		return oned.NewITFWriter()
	case QRCode:
		return qrcode.NewQRCodeWriter()
	case UPCA:
		return oned.NewUPCAWriter()
	default:
		return nil
	}
}
