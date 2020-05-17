//
// go-rencode v0.1.4 - Go implementation of rencode - fast (basic)
//                  object serialization similar to bencode
// Copyright (C) 2015~2019 gdm85 - https://github.com/gdm85/go-rencode/

// This program is free software; you can redistribute it and/or
// modify it under the terms of the GNU General Public License
// as published by the Free Software Foundation; either version 2
// of the License, or (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program; if not, write to the Free Software
// Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301, USA.

package rencode

import (
	"errors"
	"fmt"
)

// ConversionOverflow is returned when the scanned integer would overflow the destination integer
type ConversionOverflow struct {
	SourceTypeName string
	DestTypeName   string
}

func (co ConversionOverflow) Error() string {
	return fmt.Sprintf("conversion from %q to %q would overflow integer size", co.SourceTypeName, co.DestTypeName)
}

// Scan will scan the decoder data to fill in the specified target objects; if possible,
// a conversion will be performed. If targets have not pointer types or if the conversion is
// not possible, an error will be returned.
func (d *Decoder) Scan(targets ...interface{}) error {
	for i, target := range targets {
		src, err := d.DecodeNext()
		if err != nil {
			return err
		}

		err = convertAssign(src, target)
		if err != nil {
			return fmt.Errorf("scan element %d: %v", i, err)
		}
	}

	return nil
}

// Scan will scan the list to fill in the specified target objects; if possible,
// a conversion will be performed. If targets have not pointer types or if the conversion is
// not possible, an error will be returned.
// 32-bit integers larger than 16777216 will be imprecisely allowed to cast to float32.
func (l *List) Scan(targets ...interface{}) error {
	if len(targets) > l.Length() {
		return errors.New("not enough elements in list")
	}
	for i, target := range targets {
		err := convertAssign(l.values[i], target)
		if err != nil {
			return fmt.Errorf("scan element %d: %v", i, err)
		}
	}

	return nil
}

// Shift will remove the specified number of elements from the front of the List and
// return how many were removed. This will be different than n if the List is shorter.
func (l *List) Shift(n int) int {
	max := len(l.values)
	if max < n {
		n = max
	}
	l.values = l.values[n:]
	return n
}
