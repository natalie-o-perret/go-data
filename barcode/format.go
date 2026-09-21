package barcode

import "github.com/makiuchi-d/gozxing"

// Format identifies a barcode symbology.
type Format int

const (
	Aztec      Format = Format(gozxing.BarcodeFormat_AZTEC)
	Codabar    Format = Format(gozxing.BarcodeFormat_CODABAR)
	Code39     Format = Format(gozxing.BarcodeFormat_CODE_39)
	Code93     Format = Format(gozxing.BarcodeFormat_CODE_93)
	Code128    Format = Format(gozxing.BarcodeFormat_CODE_128)
	DataMatrix Format = Format(gozxing.BarcodeFormat_DATA_MATRIX)
	EAN8       Format = Format(gozxing.BarcodeFormat_EAN_8)
	EAN13      Format = Format(gozxing.BarcodeFormat_EAN_13)
	ITF        Format = Format(gozxing.BarcodeFormat_ITF)
	QRCode     Format = Format(gozxing.BarcodeFormat_QR_CODE)
	RSS14      Format = Format(gozxing.BarcodeFormat_RSS_14)
	UPCA       Format = Format(gozxing.BarcodeFormat_UPC_A)
)

// String returns the conventional name of the format.
func (f Format) String() string {
	return gozxing.BarcodeFormat(f).String()
}

// Dimension returns 1 or 2 for a known format and 0 otherwise.
func (f Format) Dimension() int {
	switch f {
	case Aztec, DataMatrix, QRCode:
		return 2
	case Codabar, Code39, Code93, Code128, EAN8, EAN13, ITF, RSS14, UPCA:
		return 1
	default:
		return 0
	}
}
