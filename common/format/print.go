package format

import (
	"fmt"
	"io"

	"github.com/beevik/etree"
	"github.com/hokaccha/go-prettyjson"
	"github.com/tliron/puccini/common/terminal"
)

var prettyjsonFormatter = prettyjson.NewFormatter()

func init() {
	prettyjsonFormatter.Indent = terminal.IndentSpaces
}

func Print(data interface{}, format string, writer io.Writer, pretty bool) error {
	// Special handling for strings
	if s, ok := data.(string); ok {
		if pretty {
			s += "\n"
		}
		_, err := fmt.Fprint(writer, s)
		return err
	}

	// Special handling for etree
	if xmlDocument, ok := data.(*etree.Document); ok {
		return PrintXMLDocument(xmlDocument, writer, pretty)
	}

	switch format {
	case "yaml", "":
		return PrintYAML(data, writer, pretty)
	case "json":
		return PrintJSON(data, writer, pretty)
	case "xml":
		return PrintXML(data, writer, pretty)
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
}

func PrintYAML(data interface{}, writer io.Writer, pretty bool) error {
	indent := "          "
	if pretty {
		indent = terminal.Indent
	}
	return WriteYAML(data, writer, indent)
}

func PrintJSON(data interface{}, writer io.Writer, pretty bool) error {
	if pretty {
		bytes, err := prettyjsonFormatter.Marshal(data)
		if err != nil {
			return err
		}
		fmt.Fprintf(writer, "%s\n", bytes)
	} else {
		return WriteJSON(data, writer, "")
	}
	return nil
}

func PrintXML(data interface{}, writer io.Writer, pretty bool) error {
	indent := ""
	if pretty {
		indent = terminal.Indent
	}
	if err := WriteXML(data, writer, indent); err != nil {
		return err
	}
	if pretty {
		fmt.Fprintln(writer)
	}
	return nil
}

func PrintXMLDocument(xmlDocument *etree.Document, writer io.Writer, pretty bool) error {
	if pretty {
		xmlDocument.Indent(terminal.IndentSpaces)
	} else {
		xmlDocument.Indent(0)
	}
	_, err := xmlDocument.WriteTo(writer)
	return err
}
