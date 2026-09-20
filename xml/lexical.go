package xml

import stdxml "encoding/xml"

func lexicalToken(token stdxml.Token) stdxml.Token {
	switch value := token.(type) {
	case stdxml.StartElement:
		value.Name = lexicalName(value.Name)
		for index := range value.Attr {
			value.Attr[index].Name = lexicalName(value.Attr[index].Name)
		}
		return value
	case stdxml.EndElement:
		value.Name = lexicalName(value.Name)
		return value
	default:
		return token
	}
}

func lexicalName(name stdxml.Name) stdxml.Name {
	if name.Space != "" {
		name.Local = name.Space + ":" + name.Local
		name.Space = ""
	}
	return name
}
