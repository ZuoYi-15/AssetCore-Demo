package data

import (
	"archive/zip"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"regexp"
	"strconv"
	"strings"

	"asset-core/internal/module/asset"
)

type AssetExcelRow struct {
	RowNumber int
	Request   asset.CreateRequest
	Raw       map[string]string
}

type sharedStringsXML struct {
	Items []struct {
		Text string `xml:"t"`
		Runs []struct {
			Text string `xml:"t"`
		} `xml:"r"`
	} `xml:"si"`
}

type worksheetXML struct {
	Rows []struct {
		Index int `xml:"r,attr"`
		Cells []struct {
			Ref    string `xml:"r,attr"`
			Type   string `xml:"t,attr"`
			Value  string `xml:"v"`
			Inline struct {
				Text string `xml:"t"`
			} `xml:"is"`
		} `xml:"c"`
	} `xml:"sheetData>row"`
}

type workbookXML struct {
	Sheets []struct {
		Name string `xml:"name,attr"`
		RID  string `xml:"http://schemas.openxmlformats.org/officeDocument/2006/relationships id,attr"`
	} `xml:"sheets>sheet"`
}

type relationshipsXML struct {
	Relationships []struct {
		ID     string `xml:"Id,attr"`
		Target string `xml:"Target,attr"`
		Type   string `xml:"Type,attr"`
	} `xml:"Relationship"`
}

var cellRefPattern = regexp.MustCompile(`^([A-Z]+)([0-9]+)$`)

func ReadAssetRows(file multipart.File) ([]AssetExcelRow, error) {
	size, err := file.Seek(0, io.SeekEnd)
	if err != nil {
		return nil, err
	}
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}
	reader, err := zip.NewReader(file, size)
	if err != nil {
		return nil, fmt.Errorf("invalid xlsx file: %w", err)
	}

	shared, err := readSharedStrings(reader)
	if err != nil {
		return nil, err
	}
	sheetPath := firstWorksheetPath(reader)
	if sheetPath == "" {
		return nil, errors.New("xlsx worksheet not found")
	}
	sheet, err := readWorksheet(reader, sheetPath)
	if err != nil {
		return nil, err
	}
	if len(sheet.Rows) < 2 {
		return nil, errors.New("xlsx has no asset rows")
	}

	header := map[int]string{}
	for _, cell := range sheet.Rows[0].Cells {
		col := columnIndex(cell.Ref)
		header[col] = normalizeHeader(cellValue(cell.Type, cell.Value, cell.Inline.Text, shared))
	}

	rows := make([]AssetExcelRow, 0, len(sheet.Rows)-1)
	for _, xrow := range sheet.Rows[1:] {
		raw := map[string]string{}
		values := map[string]string{}
		for _, cell := range xrow.Cells {
			col := columnIndex(cell.Ref)
			key := header[col]
			if key == "" {
				continue
			}
			value := strings.TrimSpace(cellValue(cell.Type, cell.Value, cell.Inline.Text, shared))
			values[key] = value
			raw[key] = value
		}
		if isEmptyRow(values) {
			continue
		}
		rows = append(rows, AssetExcelRow{
			RowNumber: xrow.Index,
			Request: asset.CreateRequest{
				AssetName:       values["asset_name"],
				AssetType:       values["asset_type"],
				Vendor:          values["vendor"],
				Model:           values["model"],
				SerialNumber:    values["serial_number"],
				MACAddress:      values["mac_address"],
				IPAddress:       values["ip_address"],
				Hostname:        values["hostname"],
				OwnerDepartment: values["owner_department"],
				OwnerUser:       values["owner_user"],
				Location:        values["location"],
				Source:          values["source"],
			},
			Raw: raw,
		})
	}
	return rows, nil
}

func readSharedStrings(reader *zip.Reader) ([]string, error) {
	file := findZipFile(reader, "xl/sharedStrings.xml")
	if file == nil {
		return nil, nil
	}
	rc, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer rc.Close()
	var doc sharedStringsXML
	if err := xml.NewDecoder(rc).Decode(&doc); err != nil {
		return nil, err
	}
	values := make([]string, 0, len(doc.Items))
	for _, item := range doc.Items {
		if item.Text != "" {
			values = append(values, item.Text)
			continue
		}
		var builder strings.Builder
		for _, run := range item.Runs {
			builder.WriteString(run.Text)
		}
		values = append(values, builder.String())
	}
	return values, nil
}

func firstWorksheetPath(reader *zip.Reader) string {
	if path := firstWorksheetPathFromWorkbook(reader); path != "" {
		return path
	}
	preferred := findZipFile(reader, "xl/worksheets/sheet1.xml")
	if preferred != nil {
		return "xl/worksheets/sheet1.xml"
	}
	for _, file := range reader.File {
		name := normalizeZipName(file.Name)
		if strings.HasPrefix(name, "xl/worksheets/") && strings.HasSuffix(strings.ToLower(name), ".xml") {
			return name
		}
	}
	return ""
}

func firstWorksheetPathFromWorkbook(reader *zip.Reader) string {
	workbookFile := findZipFile(reader, "xl/workbook.xml")
	relsFile := findZipFile(reader, "xl/_rels/workbook.xml.rels")
	if workbookFile == nil || relsFile == nil {
		return ""
	}

	workbookRC, err := workbookFile.Open()
	if err != nil {
		return ""
	}
	defer workbookRC.Close()
	var workbook workbookXML
	if err := xml.NewDecoder(workbookRC).Decode(&workbook); err != nil || len(workbook.Sheets) == 0 {
		return ""
	}

	relsRC, err := relsFile.Open()
	if err != nil {
		return ""
	}
	defer relsRC.Close()
	var rels relationshipsXML
	if err := xml.NewDecoder(relsRC).Decode(&rels); err != nil {
		return ""
	}

	firstRID := workbook.Sheets[0].RID
	for _, rel := range rels.Relationships {
		if rel.ID != firstRID {
			continue
		}
		target := normalizeZipName(rel.Target)
		if strings.HasPrefix(target, "/") {
			return strings.TrimPrefix(target, "/")
		}
		if strings.HasPrefix(target, "xl/") {
			return target
		}
		return "xl/" + target
	}
	return ""
}

func readWorksheet(reader *zip.Reader, path string) (*worksheetXML, error) {
	file := findZipFile(reader, path)
	if file == nil {
		return nil, errors.New("worksheet not found")
	}
	rc, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer rc.Close()
	var sheet worksheetXML
	if err := xml.NewDecoder(rc).Decode(&sheet); err != nil {
		return nil, err
	}
	return &sheet, nil
}

func findZipFile(reader *zip.Reader, name string) *zip.File {
	want := normalizeZipName(name)
	for _, file := range reader.File {
		if normalizeZipName(file.Name) == want {
			return file
		}
	}
	return nil
}

func normalizeZipName(name string) string {
	name = strings.ReplaceAll(name, "\\", "/")
	name = strings.TrimPrefix(name, "./")
	return strings.TrimLeft(name, "/")
}

func cellValue(cellType, value, inline string, shared []string) string {
	switch cellType {
	case "s":
		idx, err := strconv.Atoi(value)
		if err == nil && idx >= 0 && idx < len(shared) {
			return shared[idx]
		}
	case "inlineStr":
		return inline
	}
	return value
}

func columnIndex(ref string) int {
	matches := cellRefPattern.FindStringSubmatch(ref)
	if len(matches) != 3 {
		return 0
	}
	col := 0
	for _, r := range matches[1] {
		col = col*26 + int(r-'A'+1)
	}
	return col
}

func normalizeHeader(header string) string {
	key := strings.ToLower(strings.TrimSpace(header))
	key = strings.ReplaceAll(key, " ", "_")
	key = strings.ReplaceAll(key, "-", "_")
	aliases := map[string]string{
		"资产名称":       "asset_name",
		"资产名":        "asset_name",
		"名称":         "asset_name",
		"资产类型":       "asset_type",
		"类型":         "asset_type",
		"厂商":         "vendor",
		"品牌":         "vendor",
		"型号":         "model",
		"序列号":        "serial_number",
		"sn":         "serial_number",
		"mac":        "mac_address",
		"mac地址":      "mac_address",
		"ip":         "ip_address",
		"ip地址":       "ip_address",
		"主机名":        "hostname",
		"部门":         "owner_department",
		"所属部门":       "owner_department",
		"负责人":        "owner_user",
		"责任人":        "owner_user",
		"位置":         "location",
		"数据来源":       "source",
		"来源":         "source",
		"asset_name": "asset_name",
		"asset_type": "asset_type",
	}
	if value, ok := aliases[key]; ok {
		return value
	}
	return key
}

func isEmptyRow(values map[string]string) bool {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return false
		}
	}
	return true
}
