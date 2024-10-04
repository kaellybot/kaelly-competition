package constants

type Source struct {
	Name string
	Icon string
	URL  string
}

type MapType string

const (
	MapTypeNormal   MapType = "normal"
	MapTypeTactical MapType = "tactical"

	KTArenaMapTemplateURL = "https://draft.ktarena.com/images/maps/1/%v/%v.webp"

	MapCount = 50
)

func GetMapSource() Source {
	return Source{
		Name: "Krosmoz Tournaments Arena",
		Icon: "https://ktarena.com/assets/img/layout/favicon.png",
		URL:  "https://ktarena.com/",
	}
}
