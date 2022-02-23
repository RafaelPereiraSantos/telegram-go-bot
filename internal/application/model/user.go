package model

type (
	User struct {
		Name string
	}

	UserPost struct {
		Title string
		Image PostImage
	}

	PostImage struct {
		Url string
	}
)
