package main

import (
	"fmt"
	"math"
)

const (
	h_key        = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	h_base       = 20037508.34
	h_deg        = math.Pi * (30.0 / 180.0)
	max_level    = 15
	max_code_len = max_level + 2
)

func calcHexSize(level int) float64 {
	return h_base / math.Pow(3.0, float64(level+1))
}

func loc2xy(lon float64, lat float64) (x float64, y float64) {
	x = lon * h_base / 180.0
	y = math.Log(math.Tan((90.0+lat)*math.Pi/360.0)) / (math.Pi / 180.0)
	y *= h_base / 180.0
	return
}

func xy2loc(x float64, y float64) (lat float64, lon float64) {
	lon = (x / h_base) * 180.0
	lat = (y / h_base) * 180.0
	lat = 180.0 / math.Pi * (2.0*math.Atan(math.Exp(lat*math.Pi/180.0)) - math.Pi/2.0)
	return
}

type Zone struct {
	lat, lon float64
	level    int
	code     string
	x, y     int64
}

func lround(v float64) int64 {
	i := int64(v)
	f := v - float64(i)
	if f > 0.5 {
		return i + 1
	}
	return i
}

func newZone(lat float64, lon float64, level int) *Zone {
	if lat < -90 || lat > 90 || lon < -180 || lon > 180 ||
		level < 0 || level > max_level {
		return nil
	}
	h_k := math.Tan(h_deg)

	zone := new(Zone)
	zone.level = level
	level += 2
	h_size := calcHexSize(level)
	lon_grid, lat_grid := loc2xy(lon, lat)
	unit_x := 6 * h_size
	unit_y := 6 * h_size * h_k
	h_pos_x := (lon_grid + lat_grid/h_k) / unit_x
	h_pos_y := (lat_grid - h_k*lon_grid) / unit_y
	h_x_0 := int64(math.Floor(h_pos_x))
	h_y_0 := int64(math.Floor(h_pos_y))
	h_x_q := h_pos_x - float64(h_x_0)
	h_y_q := h_pos_y - float64(h_y_0)
	h_x := lround(h_pos_x)
	h_y := lround(h_pos_y)
	if h_y_q > -h_x_q+1 {
		if (h_y_q < 2*h_x_q) && (h_y_q > 0.5*h_x_q) {
			h_x = h_x_0 + 1
			h_y = h_y_0 + 1
		}
	} else if h_y_q < -h_x_q+1 {
		if (h_y_q > (2*h_x_q)-1) && (h_y_q < (0.5*h_x_q)+0.5) {
			h_x = h_x_0
			h_y = h_y_0
		}
	}

	h_lat := (h_k*float64(h_x)*unit_x + float64(h_y)*unit_y) / 2
	h_lon := (h_lat - float64(h_y)*unit_y) / h_k

	z_loc_x, z_loc_y := xy2loc(h_lon, h_lat)
	if h_base-h_lon < h_size {
		z_loc_x = 180
		h_xy := h_x
		h_x = h_y
		h_y = h_xy
	}

	var code3_x [max_level + 1]int
	var code3_y [max_level + 1]int
	mod_x := float64(h_x)
	mod_y := float64(h_y)

	for i := 0; i <= level; i++ {
		h_pow := math.Pow(3.0, float64(level-i))
		if mod_x >= math.Ceil(h_pow/2) {
			code3_x[i] = 2
			mod_x -= h_pow
		} else if mod_x <= -math.Ceil(h_pow/2) {
			code3_x[i] = 0
			mod_x += h_pow
		} else {
			code3_x[i] = 1
		}

		if mod_y >= math.Ceil(h_pow/2) {
			code3_y[i] = 2
			mod_y -= h_pow
		} else if mod_y <= -math.Ceil(h_pow/2) {
			code3_y[i] = 0
			mod_y += h_pow
		} else {
			code3_y[i] = 1
		}
	}

	h_1 := 0
	for i := 0; i < 3; i++ {
		d := code3_x[i]*3 + code3_y[i]
		h_1 = h_1*10 + d
	}
	h_a1 := h_1 / 30
	h_a2 := h_1 % 30
	code := make([]byte, level+1)
	code[0] = h_key[h_a1]
	code[1] = h_key[h_a2]

	for i := 3; i <= level; i++ {
		d := code3_x[i]*3 + code3_y[i]
		code[i] = uint8('0' + d)
	}
	zone.code = string(code)
	zone.lat = z_loc_y
	zone.lon = z_loc_x
	zone.x = h_x
	zone.y = h_y
	return zone
}

func main() {
	fmt.Println("hello")
	z := newZone(30.0, 120.1, 5)
	if z != nil {
		fmt.Println(z.code)
	}
}
