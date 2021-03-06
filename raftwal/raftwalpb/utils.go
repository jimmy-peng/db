package raftwalpb

import "errors"

var ErrCRCMismatch = errors.New("raftwalpb: crc mismatch")

// Validate returns error if crc is not equal to Record's CRC.
func (rec *Record) Validate(crc uint32) error {
	if rec.CRC == crc {
		return nil
	}
	rec.Reset()
	return ErrCRCMismatch
}
