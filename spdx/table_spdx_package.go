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

func tableSpdxPackage(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "spdx_package",
		Description: "Details of packages in SPDX files.",
		List: &plugin.ListConfig{
			Hydrate:    listSpdxPackages,
			KeyColumns: plugin.SingleColumn("directory"),
		},
		Columns: []*plugin.Column{
			{Name: "file_path", Type: proto.ColumnType_STRING, Description: "The path to the SPDX file."},
			{Name: "package_name", Type: proto.ColumnType_STRING, Description: "The package name."},
			{Name: "package_version", Type: proto.ColumnType_STRING, Description: "The package version."},
			{Name: "package_supplier", Type: proto.ColumnType_STRING, Description: "The package supplier."},
			{Name: "package_download_location", Type: proto.ColumnType_STRING, Description: "The package download location."},
			{Name: "directory", Type: proto.ColumnType_STRING, Description: "The directory containing SPDX files."},
		},
	}
}

func listSpdxPackages(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
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

			packages := parseSPDXPackages(file)
			for _, pkg := range packages {
				pkg["file_path"] = path
				pkg["directory"] = directory
				d.StreamListItem(ctx, pkg)
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return nil, nil
}

func parseSPDXPackages(file *os.File) []map[string]interface{} {
	var packages []map[string]interface{}
	var pkg map[string]interface{}
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
		case "PackageName":
			if pkg != nil {
				packages = append(packages, pkg)
			}
			pkg = make(map[string]interface{})
			pkg["package_name"] = value
		case "PackageVersion":
			pkg["package_version"] = value
		case "PackageSupplier":
			pkg["package_supplier"] = strings.Split(value, ": ")[1]
		case "PackageDownloadLocation":
			pkg["package_download_location"] = value
		}
	}

	if pkg != nil {
		packages = append(packages, pkg)
	}

	return packages
}
