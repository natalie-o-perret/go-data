package xml

import (
	"bytes"
	stdxml "encoding/xml"
	"strings"
)

func renderDocument(document parsedDocument) ([]byte, error) {
	var output bytes.Buffer
	encoder := stdxml.NewEncoder(&output)
	var stack []int
	nextElement, topLevel := 0, 0
	indent := func(depth int) error {
		return encoder.EncodeToken(stdxml.CharData("\n" + strings.Repeat("  ", depth)))
	}
	beforeMarkup := func() error {
		if len(stack) == 0 {
			if topLevel != 0 {
				if err := indent(0); err != nil {
					return err
				}
			}
			topLevel++
			return nil
		}
		parent := document.elements[stack[len(stack)-1]]
		if parent.markup && !parent.text {
			return indent(len(stack))
		}
		return nil
	}

	for _, token := range document.tokens {
		switch value := token.(type) {
		case stdxml.StartElement:
			if err := beforeMarkup(); err != nil {
				return nil, err
			}
			if err := encoder.EncodeToken(value); err != nil {
				return nil, err
			}
			stack = append(stack, nextElement)
			nextElement++
		case stdxml.EndElement:
			current := document.elements[stack[len(stack)-1]]
			if current.markup && !current.text {
				if err := indent(len(stack) - 1); err != nil {
					return nil, err
				}
			}
			if err := encoder.EncodeToken(value); err != nil {
				return nil, err
			}
			stack = stack[:len(stack)-1]
		case stdxml.CharData:
			if len(stack) == 0 {
				continue
			}
			current := document.elements[stack[len(stack)-1]]
			if current.markup && !current.text && len(bytes.TrimSpace(value)) == 0 {
				continue
			}
			if err := encoder.EncodeToken(value); err != nil {
				return nil, err
			}
		case stdxml.Comment, stdxml.ProcInst, stdxml.Directive:
			if err := beforeMarkup(); err != nil {
				return nil, err
			}
			if err := encoder.EncodeToken(token); err != nil {
				return nil, err
			}
		}
	}

	if err := encoder.Close(); err != nil {
		return nil, err
	}
	return append(output.Bytes(), '\n'), nil
}
