package spdx

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
)

type PackageInfo struct {
	PackageName             string
	PackageVersion          string
	PackageSupplier         string
	PackageDownloadLocation string
}

func tableSpdxFile(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "spdx_file",
		Description: "Details of SPDX files.",
		List: &plugin.ListConfig{
			Hydrate:    listSpdxFiles,
			KeyColumns: plugin.SingleColumn("directory"),
		},
		Columns: []*plugin.Column{
			{Name: "path", Type: proto.ColumnType_STRING, Description: "The path to the SPDX file."},
			{Name: "spdx_version", Type: proto.ColumnType_STRING, Description: "The SPDX version."},
			{Name: "data_license", Type: proto.ColumnType_STRING, Description: "The data license."},
			{Name: "spdx_identifier", Type: proto.ColumnType_STRING, Description: "The SPDX identifier."},
			{Name: "document_name", Type: proto.ColumnType_STRING, Description: "The document name."},
			{Name: "document_namespace", Type: proto.ColumnType_STRING, Description: "The document namespace."},
			{Name: "created", Type: proto.ColumnType_STRING, Description: "The creation date."},
			{Name: "creators", Type: proto.ColumnType_STRING, Description: "The creators of the document."},
			{Name: "packages", Type: proto.ColumnType_JSON, Description: "The packages in the document."},
			{Name: "directory", Type: proto.ColumnType_STRING, Description: "The directory containing SPDX files."},
		},
	}
}

func listSpdxFiles(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	quals := d.EqualsQuals
	directory := quals["directory"].GetStringValue()

	if directory == "" {
		return nil, fmt.Errorf("directory must be specified")
	}

	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(info.Name()) == ".spdx" {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			doc := parseSPDXFile(file)
			doc["path"] = path
			doc["directory"] = directory
			d.StreamListItem(ctx, doc)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return nil, nil
}

func parseSPDXFile(file *os.File) map[string]interface{} {
	doc := make(map[string]interface{})
	var pkg PackageInfo
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, ": ", 2)
		if len(parts) != 2 {
			continue
		}

		key := parts[0]
		value := parts[1]

		switch key {
		case "SPDXVersion":
			doc["spdx_version"] = value
		case "DataLicense":
			doc["data_license"] = value
		case "SPDXID":
			doc["spdx_identifier"] = value
		case "DocumentName":
			doc["document_name"] = value
		case "DocumentNamespace":
			doc["document_namespace"] = value
		case "Created":
			doc["created"] = value
		case "Creator":
			doc["creators"] = strings.Split(value, ": ")[1]
		case "PackageName":
			if pkg.PackageName != "" {
				if doc["packages"] == nil {
					doc["packages"] = []PackageInfo{}
				}
				doc["packages"] = append(doc["packages"].([]PackageInfo), pkg)
				pkg = PackageInfo{}
			}
			pkg.PackageName = value
		case "PackageVersion":
			pkg.PackageVersion = value
		case "PackageSupplier":
			pkg.PackageSupplier = value
		case "PackageDownloadLocation":
			pkg.PackageDownloadLocation = value
		}
	}

	if pkg.PackageName != "" {
		if doc["packages"] == nil {
			doc["packages"] = []PackageInfo{}
		}
		doc["packages"] = append(doc["packages"].([]PackageInfo), pkg)
	}

	return doc
}
