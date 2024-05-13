package extype

import (
	"bytes"
	"context"
	"database/sql/driver"
	"encoding/binary"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Point struct {
	Lng, Lat float64
}

func (p Point) GormDataType() string {
	return "geometry"
}

func (p Point) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	return clause.Expr{
		SQL:  "ST_PointFromText(?)",
		Vars: []any{p.String()},
	}
}

func (p *Point) String() string {
	return fmt.Sprintf("POINT(%v %v)", p.Lng, p.Lat)
}

func (p *Point) Scan(value any) error {
	data, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("got data of type %T but wanted []byte", value)
	}

	// 去掉SRID (4 bytes)
	wkb := data[4:]
	r := bytes.NewReader(wkb)
	var byteOrder binary.ByteOrder
	byteOrderMark, err := r.ReadByte()
	if err != nil {
		return fmt.Errorf("could not read byte order: %w", err)
	}

	switch byteOrderMark {
	case 0:
		byteOrder = binary.BigEndian
	case 1:
		byteOrder = binary.LittleEndian
	default:
		return fmt.Errorf("invalid byte order %d", byteOrderMark)
	}

	var geometryType uint32
	if err = binary.Read(r, byteOrder, &geometryType); err != nil {
		return fmt.Errorf("could not read geometry type: %w", err)
	}

	if geometryType != 1 {
		return fmt.Errorf("invalid geometry type %d", geometryType)
	}

	if err = binary.Read(r, byteOrder, &p.Lng); err != nil {
		return fmt.Errorf("could not read X value: %w", err)
	}

	if err = binary.Read(r, byteOrder, &p.Lat); err != nil {
		return fmt.Errorf("could not read Y value: %w", err)
	}

	return nil
}

func (p *Point) Value() (driver.Value, error) {
	return p.String(), nil
}

type NullPoint struct {
	Point Point
	Valid bool
}

func (np *NullPoint) Scan(val any) error {
	if val == nil {
		np.Point, np.Valid = Point{}, false
		return nil
	}

	point := &Point{}
	err := point.Scan(val)
	if err != nil {
		np.Point, np.Valid = Point{}, false
		return nil
	}
	np.Point = Point{
		Lng: point.Lng,
		Lat: point.Lat,
	}
	np.Valid = true

	return nil
}

func (np NullPoint) Value() (driver.Value, error) {
	if !np.Valid {
		return nil, nil
	}

	return np.Point, nil
}
