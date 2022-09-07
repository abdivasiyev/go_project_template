package helpers

import "database/sql"

func ToNullString(s string) sql.NullString {
	return sql.NullString{
		String: s,
		Valid:  s != "",
	}
}

func ToNullByte(n byte) sql.NullByte {
	return sql.NullByte{
		Byte:  n,
		Valid: n > 0.0,
	}
}

func ToNullFloat64(n float64) sql.NullFloat64 {
	return sql.NullFloat64{
		Float64: n,
		Valid:   true,
	}
}

func ToNullInt16(n int16) sql.NullInt16 {
	return sql.NullInt16{
		Int16: n,
		Valid: true,
	}
}

func ToNullInt32(n int32) sql.NullInt32 {
	return sql.NullInt32{
		Int32: n,
		Valid: true,
	}
}

func ToNullInt64(n int64) sql.NullInt64 {
	return sql.NullInt64{
		Int64: n,
		Valid: true,
	}
}
