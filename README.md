## What is this?

This tool is made for downloading Sentinel data from the Copernicus Open Access Hub. Given a list of coordinates, it will download
all the data that contains any of the coordinates. It does so by first building a database of the files that contain the coordinates
and then requesting the data from the Copernicus Open Access Hub. If the file is not online it handles the request to the
historical archives and the quota limit.

## How to use it?

Run the command line with the path to the settings file. A sample settings file is provided in the folder as `settings.yaml`

## Caveats

The platform and product types are hardcoded, so if you want to use this tool for another platform or product type you will need to
modify the code at `builder.go`