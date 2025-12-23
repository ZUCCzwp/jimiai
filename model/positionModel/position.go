package positionModel

import (
	"math"
)

type Position struct {
	Longitude float64 `json:"longitude" gorm:"type:float"` // 经度 113.4022203113
	Latitude  float64 `json:"latitude" gorm:"type:float"`  // 纬度 23.1378010917
}

// Distance 计算距离, 单位: km
func (p Position) Distance(longitude, latitude float64) float64 {
	radius := 6371000.0 //6378137.0
	rad := math.Pi / 180.0
	lat1 := latitude * rad
	lng1 := longitude * rad
	lat2 := p.Latitude * rad
	lng2 := p.Longitude * rad

	theta := lng2 - lng1
	dist := math.Acos(math.Sin(lat1)*math.Sin(lat2) + math.Cos(lat1)*math.Cos(lat2)*math.Cos(theta))
	return dist * radius / 1000
}

type Rect struct {
	Lng1 float64
	Lat1 float64
	Lng2 float64
	Lat2 float64
}

// Rect 返回以当前坐标为中心, distance 为半径的矩形区域
// 高精度库 "github.com/kellydunn/golang-geo"
func (p Position) Rect(distance float64) Rect {
	radius := 6371.01 // 地球半径（单位：km）

	// 将距离转换成弧度
	dist := distance / radius
	latRad := p.Latitude * math.Pi / 180.0
	lngRad := p.Longitude * math.Pi / 180.0

	// 计算矩形四个边角点的经纬度
	minLat := latRad - dist
	maxLat := latRad + dist
	minLng := lngRad - dist/math.Cos(latRad)
	maxLng := lngRad + dist/math.Cos(latRad)

	rect := Rect{
		Lng1: minLng * 180.0 / math.Pi,
		Lat1: minLat * 180.0 / math.Pi,
		Lng2: maxLng * 180.0 / math.Pi,
		Lat2: maxLat * 180.0 / math.Pi,
	}
	return rect
}
