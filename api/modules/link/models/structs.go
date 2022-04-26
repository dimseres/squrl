package models

type AddLinkRequest struct {
	Url string `form:"url" json:"url" binding:"required"`
}
