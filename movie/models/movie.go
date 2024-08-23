package models

type Movie struct {
	Title         string   `bson:"title"`
	Description   string   `bson:"description"`
	Year          int      `bson:"year"`
	Duration      string   `bson:"duration"`
	Genre         []string `bson:"genre"`
	CoverUrl      string   `bson:"coverUrl"`
	BackgroundUrl string   `bson:"backgroundUrl"`
	LogoUrl       string   `bson:"logoUrl"`
}
