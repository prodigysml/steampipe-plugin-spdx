install:
	mkdir -p ~/.steampipe/plugins/local/spdx/
	go build -o  ~/.steampipe/plugins/local/spdx/spdx.plugin *.go
	mkdir -p ~/.steampipe/config/
	cp spdx.spc ~/.steampipe/config/