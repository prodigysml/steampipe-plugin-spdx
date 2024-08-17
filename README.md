# Steampipe Plugin spdx

This repository contains a steampipe plugin that parses spdx files and allows for easier searching.

## Requirements

* Golang (tested on 1.22.5)
* Steampipe

## Installation

```bash
task install
```

## Example Query

```sql
select
    f.path as spdx_file_path,
    f.document_name,
    p.package_name,
    p.package_version,
    p.package_supplier,
    p.package_download_location
from
    spdx_file f
join
    spdx_package p
on
    f.path = p.file_path
where
    f.directory = './examples' and
    p.directory = './examples';
```

## Test

```bash
task test
```

## Uninstall

```bash
task clean
```