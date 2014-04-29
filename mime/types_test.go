package MIME

import(
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestByMIMEType(t *testing.T) {
	byMIMEType := ByMIMEType()
	assert.NotNil(t, byMIMEType)
	assert.Equal(t, byMIMEType["image/gif"][0], "gif")
	assert.Nil(t, byMIMEType["unknown/type"])
}

func TestByExtension(t *testing.T) {
	byExtension := ByExtension()
	assert.NotNil(t, byExtension)
	assert.Equal(t, byExtension["gif"][0], "image/gif")
	assert.Nil(t, byExtension["unknown"])
}

func TestAddType(t *testing.T) {
	byExtension := ByExtension()
	assert.NotNil(t, byExtension)
	assert.Nil(t, byExtension["foo"])
	assert.Nil(t, byExtension["bar"])

	byMIMEType := ByMIMEType()
	assert.NotNil(t, byMIMEType)
	assert.Nil(t, byMIMEType["foo/foo"])
	assert.Nil(t, byMIMEType["foo/bar"])

	AddType([]string{"foo/foo"}, []string{"foo"})
	assert.Equal(t, byMIMEType["foo/foo"][0], "foo")
	assert.Equal(t, byExtension["foo"][0], "foo/foo")
	assert.Nil(t, byExtension["bar"])
	assert.Nil(t, byMIMEType["foo/bar"])

	AddType([]string{"foo/bar"}, []string{"bar"})
	assert.Equal(t, byMIMEType["foo/foo"][0], "foo")
	assert.Equal(t, byExtension["foo"][0], "foo/foo")
	assert.Equal(t, byMIMEType["foo/bar"][0], "bar")
	assert.Equal(t, byExtension["bar"][0], "foo/bar")

	AddType([]string{"foo/foo"}, []string{"bar"})
	assert.Equal(t, byMIMEType["foo/foo"], []string{"foo", "bar"})
	assert.Equal(t, byExtension["foo"][0], "foo/foo")

	AddType([]string{"foo/bar"}, []string{"foo"})
	assert.Equal(t, byMIMEType["foo/foo"], []string{"foo", "bar"})
	assert.Equal(t, byExtension["foo"], []string{"foo/foo", "foo/bar"})
	assert.Equal(t, byMIMEType["foo/bar"], []string{"bar", "foo"})
	assert.Equal(t, byExtension["bar"], []string{"foo/bar", "foo/foo"})
}

func TestTypeFromExtension(t *testing.T) {
	assert.Equal(t, TypeFromExtension("gif")[0], "image/gif")
	assert.Equal(t, TypeFromExtension("jpg")[0], "image/jpeg")
	assert.Equal(t, TypeFromExtension("ico")[0], "image/x-icon")
	assert.Equal(t, TypeFromExtension("unknown"), []string{})
}

func TestExtensionFromType(t *testing.T) {
	assert.Equal(t, ExtensionFromType("image/gif")[0], "gif")
	assert.Equal(t, ExtensionFromType("image/jpeg"), []string{"jpeg","jpg"})
	assert.Equal(t, ExtensionFromType("image/x-icon")[0], "ico")
	assert.Equal(t, ExtensionFromType("unknown/type"), []string{})
}

func TestTypeFromFilename(t *testing.T) {
	assert.Equal(t, TypeFromFilename("foo.bar"), []string{"foo/bar", "foo/foo"})
	assert.Equal(t, TypeFromFilename("foo.htm"), []string{"text/html"})
	assert.Equal(t, TypeFromFilename("foo.pdf"), []string{"application/pdf"})

	assert.Equal(t, TypeFromFilename("/bar/foo.bar"), []string{"foo/bar", "foo/foo"})
	assert.Equal(t, TypeFromFilename("/bar/foo.htm"), []string{"text/html"})
	assert.Equal(t, TypeFromFilename("/bar/foo.pdf"), []string{"application/pdf"})

	assert.Equal(t, TypeFromFilename("//foo.com/bar/foo.bar"), []string{"foo/bar", "foo/foo"})
	assert.Equal(t, TypeFromFilename("//foo.com/bar/foo.htm"), []string{"text/html"})
	assert.Equal(t, TypeFromFilename("//foo.com/bar/foo.pdf"), []string{"application/pdf"})

	assert.Equal(t, TypeFromFilename("http://foo.com/bar/foo.bar"), []string{"foo/bar", "foo/foo"})
	assert.Equal(t, TypeFromFilename("http://foo.com/bar/foo.htm"), []string{"text/html"})
	assert.Equal(t, TypeFromFilename("http://foo.com/bar/foo.pdf"), []string{"application/pdf"})

	assert.Equal(t, TypeFromFilename("C:\\foo.com\\bar\\foo.bar"), []string{"foo/bar", "foo/foo"})
	assert.Equal(t, TypeFromFilename("C:\\foo.com\\bar\\foo.htm"), []string{"text/html"})
	assert.Equal(t, TypeFromFilename("C:\\foo.com\\bar\\foo.pdf"), []string{"application/pdf"})
}

func TestDetectFromContentType(t *testing.T) {
	assert.Equal(t, DetectFromContentType("text/plain;charset=UTF-8"), []string{"txt"})
	assert.Equal(t, DetectFromContentType("text/plain"), []string{"txt"})
	assert.Equal(t, DetectFromContentType("image/png"), []string{"png"})
	assert.Equal(t, DetectFromContentType("image/jpeg"), []string{"jpeg", "jpg"})
}
