package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
)

func bindata_read(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	return buf.Bytes(), nil
}

func assets_demo_app_readme_md() ([]byte, error) {
	return bindata_read([]byte{
		0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x00, 0xff, 0x6c, 0x8d,
		0x31, 0x0b, 0xc2, 0x30, 0x10, 0x85, 0xf7, 0xfb, 0x15, 0x07, 0xdd, 0x44,
		0xda, 0x2e, 0x22, 0x14, 0x1c, 0x9c, 0x3a, 0x89, 0xa0, 0x93, 0x5b, 0xcf,
		0x1a, 0x92, 0x40, 0xc8, 0x1d, 0xc9, 0x15, 0xe9, 0xbf, 0x37, 0xa1, 0x8b,
		0xa0, 0x37, 0xde, 0xfb, 0xde, 0xfb, 0x1e, 0xbc, 0x24, 0x1c, 0x59, 0x67,
		0x47, 0x48, 0x22, 0xc1, 0xcf, 0xa4, 0x9e, 0x23, 0x9c, 0xfe, 0x1f, 0x40,
		0xd3, 0xe0, 0x68, 0x54, 0x7d, 0xb4, 0x98, 0x95, 0x92, 0x9a, 0x17, 0xc0,
		0x0e, 0x6f, 0x4b, 0xc4, 0x69, 0x9a, 0xec, 0x36, 0xe4, 0x63, 0x89, 0x42,
		0x28, 0x8f, 0x12, 0xdd, 0x2b, 0x85, 0x6b, 0xf5, 0x7c, 0x09, 0xf6, 0x68,
		0x5a, 0xdb, 0xd6, 0xce, 0x65, 0x3d, 0x8b, 0x6c, 0xe8, 0x55, 0x4c, 0xfc,
		0x21, 0xcb, 0x1a, 0x12, 0x3e, 0x13, 0xbf, 0xb3, 0x49, 0x03, 0x3a, 0x55,
		0x19, 0xba, 0x2e, 0xf0, 0x4c, 0xc1, 0x71, 0xd6, 0xe1, 0xd8, 0x1f, 0x7a,
		0xf8, 0x04, 0x00, 0x00, 0xff, 0xff, 0x72, 0xc0, 0xcf, 0x9f, 0xc7, 0x00,
		0x00, 0x00,
	},
		"assets/demo_app/README.md",
	)
}

func assets_demo_app_assets_templates_index_html() ([]byte, error) {
	return bindata_read([]byte{
		0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x00, 0xff, 0x2c, 0x8e,
		0x3d, 0x0e, 0xc2, 0x30, 0x0c, 0x85, 0xe7, 0x72, 0x0a, 0x33, 0x21, 0x16,
		0xaa, 0xae, 0xc8, 0x0a, 0x23, 0x17, 0x40, 0xea, 0x6c, 0x12, 0x4b, 0x89,
		0x94, 0xd6, 0x11, 0x35, 0x43, 0x54, 0xf5, 0xee, 0xe4, 0x87, 0xc5, 0xcf,
		0xb2, 0xbf, 0x67, 0x3f, 0x3c, 0x3b, 0xb1, 0x9a, 0x13, 0x83, 0xd7, 0x25,
		0x9a, 0x13, 0x76, 0x19, 0xd0, 0x33, 0xb9, 0xa2, 0x03, 0x6a, 0xd0, 0xc8,
		0x66, 0xdf, 0x6f, 0xaf, 0xda, 0x1c, 0x07, 0x8e, 0x7d, 0x52, 0x98, 0xf1,
		0x0f, 0xe1, 0x5b, 0x5c, 0x6e, 0xb0, 0x9f, 0xcc, 0xcc, 0xd1, 0xca, 0xc2,
		0xa0, 0x02, 0x59, 0xbe, 0x1f, 0x78, 0x8a, 0x5a, 0x4f, 0x40, 0x29, 0xc5,
		0x60, 0x49, 0x83, 0xac, 0xc5, 0x37, 0x35, 0x3a, 0xd5, 0x3a, 0xcc, 0xb4,
		0xad, 0x17, 0x05, 0xf5, 0xa4, 0xc0, 0xb4, 0xe5, 0x07, 0xdc, 0xaf, 0x75,
		0x3b, 0xa6, 0xf6, 0xa3, 0xdf, 0x2e, 0x9e, 0x16, 0xec, 0x17, 0x00, 0x00,
		0xff, 0xff, 0x85, 0x16, 0x11, 0xcc, 0xb0, 0x00, 0x00, 0x00,
	},
		"assets/demo_app/assets/templates/index.html",
	)
}

func assets_demo_app_makefile() ([]byte, error) {
	return bindata_read([]byte{
		0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x00, 0xff, 0x74, 0x90,
		0xb1, 0x4e, 0x33, 0x31, 0x10, 0x84, 0xeb, 0xdf, 0x4f, 0x31, 0xc5, 0x2f,
		0x05, 0x8a, 0xf3, 0x89, 0xf6, 0x24, 0x3a, 0x28, 0x68, 0x42, 0x24, 0x68,
		0x28, 0x9d, 0xdc, 0xc6, 0x67, 0xe4, 0xf3, 0x45, 0xde, 0x8d, 0x38, 0x64,
		0xfc, 0xee, 0xd8, 0x39, 0x92, 0x54, 0x69, 0x2c, 0x6b, 0x77, 0xbe, 0x99,
		0xb1, 0x9f, 0x9e, 0x37, 0x6f, 0x78, 0xc4, 0xff, 0x3b, 0x3b, 0xc1, 0x3b,
		0x16, 0x34, 0x7b, 0xac, 0x52, 0x8a, 0x26, 0x58, 0x82, 0x7e, 0x27, 0x96,
		0x97, 0xf1, 0x30, 0x45, 0xe1, 0x9c, 0x53, 0xd2, 0x39, 0x23, 0x25, 0x0a,
		0x7d, 0xce, 0x2b, 0xe8, 0x56, 0x6b, 0x7d, 0xaf, 0x94, 0xf1, 0xbe, 0x43,
		0x4f, 0x07, 0x56, 0xff, 0xec, 0xd4, 0x6c, 0x5d, 0xe8, 0x8d, 0x18, 0x18,
		0x66, 0x12, 0xae, 0x92, 0x3a, 0x86, 0x0b, 0x2c, 0x45, 0xb8, 0x40, 0x4a,
		0x49, 0xf1, 0xed, 0x50, 0xcf, 0xe6, 0x4c, 0x2e, 0xe9, 0xa7, 0x3d, 0x7e,
		0x30, 0x9b, 0x68, 0x19, 0x4d, 0x78, 0x40, 0xd9, 0x54, 0x9d, 0x52, 0x91,
		0x3c, 0x19, 0xa6, 0x0e, 0x7f, 0x97, 0x0b, 0x39, 0x9f, 0x5d, 0xeb, 0xa0,
		0x3b, 0x79, 0x59, 0x12, 0x58, 0x27, 0xc3, 0x71, 0xab, 0x77, 0xd3, 0xd8,
		0x7e, 0x0a, 0xd1, 0xf1, 0x8b, 0x42, 0x7b, 0x2d, 0x78, 0xed, 0xd1, 0xdc,
		0xc4, 0x58, 0x22, 0xc9, 0x6e, 0x88, 0x6d, 0x95, 0xb9, 0xfd, 0xf7, 0xa5,
		0xc3, 0x6d, 0x64, 0x74, 0x05, 0x20, 0xef, 0x87, 0x12, 0x35, 0x2b, 0xa5,
		0x37, 0xeb, 0xd7, 0xf5, 0x47, 0x87, 0xfa, 0xf4, 0x6a, 0xb2, 0xfc, 0xd3,
		0x6f, 0x00, 0x00, 0x00, 0xff, 0xff, 0x28, 0x56, 0xdc, 0x19, 0x74, 0x01,
		0x00, 0x00,
	},
		"assets/demo_app/Makefile",
	)
}

func assets_demo_app_main_go2() ([]byte, error) {
	return bindata_read([]byte{
		0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x00, 0xff, 0x94, 0x91,
		0xc1, 0x8b, 0xea, 0x30, 0x10, 0xc6, 0xcf, 0xcd, 0x5f, 0x11, 0x72, 0x4a,
		0x45, 0x53, 0xf0, 0xf8, 0xde, 0xf3, 0xf0, 0x10, 0x9e, 0xef, 0xb8, 0xa8,
		0xe0, 0x61, 0xd9, 0x43, 0xac, 0xd3, 0x36, 0x98, 0x26, 0x25, 0x19, 0x17,
		0x61, 0xf1, 0x7f, 0xdf, 0x49, 0x5b, 0x51, 0x16, 0xd9, 0xc5, 0x53, 0x99,
		0xe9, 0x37, 0x5f, 0x7e, 0xdf, 0x4c, 0xa7, 0xcb, 0xa3, 0xae, 0x81, 0xb7,
		0xda, 0x38, 0xc6, 0x4c, 0xdb, 0xf9, 0x80, 0x92, 0x65, 0xc2, 0xfa, 0x5a,
		0xb0, 0xac, 0xf6, 0x58, 0x36, 0x9a, 0x8b, 0xda, 0x60, 0x73, 0xda, 0xab,
		0xd2, 0xb7, 0x85, 0xd1, 0x6e, 0x76, 0x04, 0x87, 0xc5, 0xf0, 0xaf, 0xd0,
		0x5d, 0x47, 0xc2, 0xef, 0x14, 0x0d, 0xe2, 0x4f, 0x92, 0xe0, 0x4f, 0x08,
		0x41, 0xb0, 0x9c, 0xb1, 0xea, 0xe4, 0xca, 0x9e, 0x46, 0xe6, 0xfc, 0x83,
		0x65, 0xef, 0x3a, 0x70, 0x7a, 0x82, 0x2f, 0xf8, 0x20, 0x55, 0xcb, 0x00,
		0x1a, 0x41, 0xfe, 0x8d, 0x11, 0x90, 0xe4, 0x59, 0xe0, 0xbf, 0x16, 0x49,
		0xa1, 0xd6, 0xbd, 0x07, 0x35, 0xd4, 0x0a, 0x50, 0x8a, 0x42, 0x4c, 0x39,
		0x9c, 0x75, 0xdb, 0x59, 0xc8, 0x6f, 0xcd, 0xca, 0xfb, 0x5b, 0x7f, 0x9e,
		0xe6, 0x29, 0xa8, 0x7a, 0x09, 0xc6, 0xa1, 0x75, 0x52, 0x6c, 0x50, 0x07,
		0x34, 0xae, 0x4e, 0x86, 0xd6, 0x94, 0x1a, 0x8d, 0x77, 0x82, 0xc6, 0x93,
		0x7f, 0xff, 0x4f, 0xa6, 0x91, 0x3f, 0xb3, 0x54, 0x2f, 0x1b, 0x76, 0x19,
		0x71, 0x47, 0x3f, 0x19, 0x21, 0x46, 0x9a, 0xe0, 0x93, 0x14, 0x59, 0x6d,
		0x86, 0x6a, 0xca, 0xfb, 0x74, 0x7c, 0x32, 0x00, 0x0e, 0x9c, 0x7d, 0xb6,
		0x51, 0x9e, 0x9c, 0x63, 0xf3, 0x2a, 0xb6, 0x06, 0x2d, 0x88, 0x37, 0x8a,
		0x2a, 0x76, 0x60, 0x69, 0x4d, 0xc0, 0xd1, 0xf3, 0x55, 0x1f, 0x5b, 0xdc,
		0xd4, 0x6b, 0x70, 0x07, 0x08, 0x52, 0x18, 0xfa, 0x9c, 0x55, 0x83, 0xad,
		0x25, 0xc2, 0x2b, 0x09, 0xe5, 0x7b, 0x82, 0x62, 0xca, 0x2b, 0x3e, 0xd6,
		0xff, 0xb5, 0x3b, 0x58, 0x08, 0xff, 0xc8, 0xe4, 0x11, 0x5b, 0xda, 0x5b,
		0x4f, 0xb6, 0xd7, 0x74, 0xa6, 0xac, 0xba, 0x3e, 0x32, 0xda, 0xe6, 0x5f,
		0x57, 0x31, 0x7f, 0x76, 0x17, 0x77, 0xe0, 0xa3, 0x88, 0xe0, 0xc8, 0xf0,
		0x19, 0x1f, 0xb2, 0xb9, 0xdb, 0x52, 0xec, 0xbc, 0x8b, 0xa0, 0x76, 0xc1,
		0x20, 0x6c, 0xe1, 0x8c, 0xf2, 0x51, 0x24, 0x25, 0x23, 0xd2, 0xf1, 0xeb,
		0x9c, 0x8e, 0x7c, 0xc9, 0x7f, 0x53, 0x8c, 0xcf, 0x00, 0x00, 0x00, 0xff,
		0xff, 0xcf, 0xeb, 0x87, 0x51, 0x11, 0x03, 0x00, 0x00,
	},
		"assets/demo_app/main.go2",
	)
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	if f, ok := _bindata[name]; ok {
		return f()
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() ([]byte, error){
	"assets/demo_app/README.md":                   assets_demo_app_readme_md,
	"assets/demo_app/assets/templates/index.html": assets_demo_app_assets_templates_index_html,
	"assets/demo_app/Makefile":                    assets_demo_app_makefile,
	"assets/demo_app/main.go2":                    assets_demo_app_main_go2,
}
