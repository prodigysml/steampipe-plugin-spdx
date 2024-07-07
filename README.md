# steampipe-plugin-spdx

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
    f.directory = './' and
    p.directory = './';
```

