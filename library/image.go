package library

type Image interface {
	GetThumbnail() string
	GetFullPath() string
	GenerateThumbnail()
  Clone() Image
}
