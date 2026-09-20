package xml

import (
	"bytes"
	stdxml "encoding/xml"
	"io"
)

type elementFormat struct {
	markup bool
	text   bool
}

type parsedDocument struct {
	tokens   []stdxml.Token
	elements []elementFormat
}

func parseDocument(data []byte) (parsedDocument, error) {
	decoder := stdxml.NewDecoder(bytes.NewReader(data))
	validator := stdxml.NewEncoder(io.Discard)
	var document parsedDocument
	var stack []int
	roots := 0

	for {
		token, err := decoder.RawToken()
		if err == io.EOF {
			break
		}
		if err != nil {
			return parsedDocument{}, err
		}
		if instruction, ok := token.(stdxml.ProcInst); ok && instruction.Target == "xml" && len(document.tokens) != 0 {
			return parsedDocument{}, errDeclaration
		}

		token = lexicalToken(token)
		if err := validator.EncodeToken(token); err != nil {
			return parsedDocument{}, err
		}

		switch value := token.(type) {
		case stdxml.StartElement:
			if len(stack) == 0 {
				roots++
				if roots > 1 {
					return parsedDocument{}, errRoot
				}
			} else {
				document.elements[stack[len(stack)-1]].markup = true
			}
			document.elements = append(document.elements, elementFormat{})
			stack = append(stack, len(document.elements)-1)
		case stdxml.EndElement:
			stack = stack[:len(stack)-1]
		case stdxml.CharData:
			if len(stack) == 0 {
				if len(bytes.TrimSpace(value)) != 0 {
					return parsedDocument{}, errText
				}
			} else if len(bytes.TrimSpace(value)) != 0 {
				document.elements[stack[len(stack)-1]].text = true
			}
		case stdxml.Comment, stdxml.ProcInst, stdxml.Directive:
			if len(stack) != 0 {
				document.elements[stack[len(stack)-1]].markup = true
			}
		}
		document.tokens = append(document.tokens, stdxml.CopyToken(token))
	}

	if roots != 1 {
		return parsedDocument{}, errRoot
	}
	if err := validator.Close(); err != nil {
		return parsedDocument{}, err
	}
	return document, nil
}
