package photolibrary

type Image struct {
	Thumbnail string
	FullPath  string
}

func NewImage(thumbnail, fullpath string) *Image {
	return &Image{Thumbnail: thumbnail, FullPath: fullpath}
}

func CopyImage(image Image) *Image {
	return &Image{Thumbnail: image.Thumbnail, FullPath: image.FullPath}
}
