package photolibrary

type Image interface {
	GetThumbnail() string
	GetFullPath() string
	GenerateThumbnail()
  Clone() Image
}
