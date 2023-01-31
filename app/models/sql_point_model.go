package models

import (
	"database/sql/driver"
	"encoding/binary"
	"fmt"
	"github.com/twpayne/go-geom/encoding/wkb"
)

type MysqlPoint struct {
	Point wkb.Point
}

func (m *MysqlPoint) Value() (driver.Value, error) {
	value, err := m.Point.Value()
	if err != nil {
		return nil, err
	}

	buf, ok := value.([]byte)
	if !ok {
		return nil, fmt.Errorf("did not convert value: expected []byte, but was %T", value)
	}

	mysqlEncoding := make([]byte, 4)
	binary.LittleEndian.PutUint32(mysqlEncoding, 4326)
	mysqlEncoding = append(mysqlEncoding, buf...)

	return mysqlEncoding, err
}

func (m *MysqlPoint) Scan(src interface{}) error {
	if src == nil {
		return nil
	}

	mysqlEncoding, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("did not scan: expected []byte but was %T", src)
	}

	var srid uint32 = binary.LittleEndian.Uint32(mysqlEncoding[0:4])

	err := m.Point.Scan(mysqlEncoding[4:])

	m.Point.SetSRID(int(srid))

	return err
}
