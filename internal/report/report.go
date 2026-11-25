package report

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"strconv"
)

type Record struct {
	SKU        string
	Qty        int
	PriceCents int // unit price
}
type jsonRow struct {
	SKU        string `json:"sku"`
	Qty        int    `json:"qty"`
	PriceCents int    `json:"price_cents"`
	LineCents  int    `json:"line_cents"`
}

type Exporter interface {
	Name() string
	Export(records []Record) ([]byte, error)
}

type JSONExporter struct{}

type CSVExporter struct {
	Comma rune // if 0, use ','
}

func (JSONExporter) Name() string {
	return "json"
}
func (JSONExporter) Export(records []Record) ([]byte, error) {
	rows := make([]jsonRow, 0, len(records))
	for _, r := range records {
		rows = append(rows, jsonRow{
			SKU:        r.SKU,
			Qty:        r.Qty,
			PriceCents: r.PriceCents,
			LineCents:  r.Qty * r.PriceCents,
		})
	}
	return json.Marshal(rows)
}

func (CSVExporter) Name() string {
	return "csv"
}
func (CSVExporter) Export(records []Record) ([]byte, error) {
	var buf bytes.Buffer
	w := csv.NewWriter(&buf)
	w.Comma = e.comma()

	if err := w.Write([]string{"sku", "qty", "price_cents", "line_cents"}); err != nil {
		return nil, err
	}

	for _, r := range records {
		row := []string{
			r.SKU,
			strconv.Itoa(r.Qty),
			strconv.Itoa(r.PriceCents),
			strconv.Itoa(r.Qty * r.PriceCents),
		}
		if err := w.Write(row); err != nil {
			return nil, err
		}
	}
	w.Flush()
	if err := w.Error(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
func (e CSVExporter) comma() rune {
	if e.Comma == 0 {
		return ','
	}
	return e.Comma
}

// RenderAll runs every exporter. Keys are exporter names.
func RenderAll(records []Record, exporters []Exporter) (map[string][]byte, error) {
	out := make(map[string][]byte, len(exporters))
	for i, e := range exporters {
		if e == nil {
			return nil, fmt.Errorf("render all: exporter #%d is nil", i)
		}
		name := e.Name()
		if _, dup := out[name]; dup {
			return nil, fmt.Errorf("render all: duplicate exporter name %q", name)
		}
		data, err := e.Export(records)
		if err != nil {
			return nil, fmt.Errorf("render all: %s: %w", name, err)
		}
		out[name] = data
	}
	return out, nil
}

func CommaOf(e Exporter) (rune, bool) {
	switch v := e.(type) {
	case CSVExporter:
		return v.comma(), true
	case *CSVExporter:
		if v == nil {
			return 0, false
		}
		return v.comma(), true
	default:
		return 0, false
	}
}
