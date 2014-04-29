package MIME

import(
	"path/filepath"
	"regexp"
	"strings"
)

var byExtension = map[string][]string{
	"appcache": []string{"text/cache-manifest"},
	"atom"    : []string{"application/atom+xml"},
	"bin"     : []string{"application/octet-stream"},
	"css"     : []string{"text/css"},
	"gif"     : []string{"image/gif"},
	"gz"      : []string{"application/x-gzip"},
	"htm"     : []string{"text/html"},
	"html"    : []string{"text/html;charset=UTF-8"},
	"ico"     : []string{"image/x-icon"},
	"jpeg"    : []string{"image/jpeg"},
	"jpg"     : []string{"image/jpeg"},
	"js"      : []string{"application/javascript"},
	"json"    : []string{"application/json"},
	"mp3"     : []string{"audio/mpeg"},
	"mp4"     : []string{"video/mp4"},
	"ogg"     : []string{"audio/ogg"},
	"ogv"     : []string{"video/ogg"},
	"pdf"     : []string{"application/pdf"},
	"png"     : []string{"image/png"},
	"rss"     : []string{"application/rss+xml"},
	"svg"     : []string{"image/svg+xml"},
	"txt"     : []string{"text/plain;charset=UTF-8", "text/plain"},
	"webm"    : []string{"video/webm"},
	"woff"    : []string{"application/font-woff"},
	"xml"     : []string{"application/xml","text/xml"},
	"zip"     : []string{"application/zip"},
};

// FIXME ordering changes because of map

var byMIMEType = make(map[string][]string, 0)
func ByMIMEType() map[string][]string {
	if len(byMIMEType) == 0 {
		for ext, mimes := range byExtension {
			for _, mime := range mimes {
				exts, ok := byMIMEType[mime]
				if ok {
					byMIMEType[mime] = append(exts, ext)
				} else {
					byMIMEType[mime] = []string{ext}
				}
			}
		}
	}
	return byMIMEType
}

func ByExtension() map[string][]string {
	return byExtension
}

func AddType(mimes []string, extensions []string) {
	for _, ext := range extensions {
		m, ok := byExtension[ext]
		if ok {
			byExtension[ext] = append(m, mimes...)
		} else {
			byExtension[ext] = mimes
		}
	}
	for _, mime := range mimes {
		exts, ok := byMIMEType[mime]
		if ok {
			byMIMEType[mime] = append(exts, extensions...)
		} else {
			byMIMEType[mime] = extensions
		}
	}
}

func TypeFromExtension(extension string) []string {
	mimes, ok := ByExtension()[extension]
	if ok { return mimes }
	return []string{}
}

func ExtensionFromType(mime string) []string {
	exts, ok := ByMIMEType()[mime]
	if ok { return exts }
	return []string{}
}

func TypeFromFilename(filename string) []string {
	ext := strings.TrimLeft(filepath.Ext(filename), ".")
	return TypeFromExtension(ext)
}

func DetectFromContentType(header string) []string {
	re := regexp.MustCompile("([^;]+)(?:;?(.*))")
	match := re.FindStringSubmatch(header)
	if len(match) == 2 {
		return ExtensionFromType(match[1])
	} else if len(match) == 3 {
		ext := ExtensionFromType(match[1] + ";" + match[2])
		if len(ext) == 0 {
			ext = ExtensionFromType(match[1])
		}
		return ext
	}
	return []string{}
}
