package geohex3

import (
	"math"
	"strings"
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
	Lat, Lon float64
	Level    int
	Code     string
	X, Y     int64
}

func lround(v float64) int64 {
	i := int64(v)
	f := v - float64(i)
	if f > 0.5 {
		return i + 1
	}
	return i
}

func GetZoneByLocation(lat float64, lon float64, level int) *Zone {
	if lat < -90 || lat > 90 || lon < -180 || lon > 180 ||
		level < 0 || level > max_level {
		return nil
	}
	h_k := math.Tan(h_deg)

	zone := new(Zone)
	zone.Level = level
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

	z_loc_y, z_loc_x := xy2loc(h_lon, h_lat)
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
	code := make([]byte, level)
	code[0] = h_key[h_a1]
	code[1] = h_key[h_a2]

	for i := 3; i <= level; i++ {
		d := code3_x[i]*3 + code3_y[i]
		code[i - 1] = uint8('0' + d)
	}
	zone.Code = string(code)
	zone.Lat = z_loc_y
	zone.Lon = z_loc_x
	zone.X = h_x
	zone.Y = h_y
	return zone
}

func GetZoneByCode(code string) *Zone {
	level := len(code);
	if (level - 2 < 0 || level - 2 > max_level) {
		return nil;
	}
	h_k := math.Tan(h_deg)
	h_size := calcHexSize(level);
	unit_x := 6.0 * h_size;
	unit_y := 6.0 * h_size * h_k;

	h_a := 0;
	cp := strings.IndexRune(h_key, int(code[0]))
	if (cp < 0) {
		return nil;
	}
	h_a += cp * 30;
	cp = strings.IndexRune(h_key, int(code[1]))
	if (cp < 0) {
		return nil;
	}
	h_a += cp

	h_a0 := (h_a % 1000) / 100;
	h_a1 := (h_a % 100) / 10;
	h_a2 := (h_a % 10);
	if (h_a1 != 1 && h_a1 != 2 && h_a1 != 5 &&
		h_a2 != 1 && h_a2 != 2 && h_a2 != 5) {
		if (h_a0 == 5) {
			h_a0 = 7;
		} else if (h_a0 == 1) {
			h_a0 = 3;
		}
	}

	h_decx := make([]int, level + 1)
	h_decy := make([]int, level + 1)

	h_decx[0] = h_a0 / 3;
	h_decy[0] = h_a0 % 3;
	h_decx[1] = h_a1 / 3;
	h_decy[1] = h_a1 % 3;
	h_decx[2] = h_a2 / 3;
	h_decy[2] = h_a2 % 3;

	for i := 3; i <= level; i++ {
		n := code[i - 1] - '0';
		if (n < 0 || n > 8) {
			return nil;
		}
		h_decx[i] = int(n / 3);
		h_decy[i] = int(n % 3);
	}

	h_x := int64(0);
	h_y := int64(0);
	for i := 0; i <= level; i++ {
		h_pow := math.Pow(3.0, float64(level - i))
		if (h_decx[i] == 0) {
			h_x -= int64(h_pow)  // XXX
		} else if (h_decx[i] == 2) {
			h_x += int64(h_pow)
		}
		if (h_decy[i] == 0) {
			h_y -= int64(h_pow)
		} else if (h_decy[i] == 2) {
			h_y += int64(h_pow)
		}
	}

	h_lat_y := (h_k * float64(h_x) * unit_x + float64(h_y) * unit_y) / 2;
	h_lon_x := (h_lat_y - float64(h_y) * unit_y) / h_k;
	h_lat, h_lon := xy2loc(h_lon_x, h_lat_y);
	if (h_lon > 180) {
		h_lon -= 360;
	} else if (h_lon < -180) {
		h_lon += 360;
	}

	zone := new(Zone)
	zone.Code = code
	zone.Lat = h_lat;
	zone.Lon = h_lon;
	zone.Level = level - 2;
	zone.X = h_x;
	zone.Y = h_y;
	return zone;
}
