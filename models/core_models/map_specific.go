package coremodels

type Address struct {
	Lat     float32 `json:"lat" binding:"required"`
	Lon     float32 `json:"lon" binding:"required"`
	Address string  `json:"address" binding:"required"`
}
