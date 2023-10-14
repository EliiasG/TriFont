package trifont

import (
	"encoding/binary"
	"io"
	"sort"
)

type Char struct {
	Vertices [][2]float32
	Indices  []uint16
	Advance  float32
}

type Font struct {
	Chars map[rune]Char
}

func (f *Font) ToBinary(w io.Writer) error {
	runes := f.getRunes()
	// header
	err := writeHeader(w, runes)
	if err != nil {
		return err
	}
	// chars
	for _, char := range f.getChars(runes) {
		err = writeChar(w, char)
		if err != nil {
			return err
		}
	}
	return nil
}

// uses int because then i do not have to implemenet the sorter
func (f *Font) getRunes() []int {
	// sort
	runes := make([]int, 0, len(f.Chars))
	for k := range f.Chars {
		runes = append(runes, int(k))
	}
	sort.Ints(runes)
	return runes
}

func (f *Font) getChars(runes []int) []Char {
	chars := make([]Char, len(runes))
	for i, run := range runes {
		chars[i] = f.Chars[rune(run)]
	}
	return chars
}

func writeHeader(w io.Writer, runes []int) error {
	err := binary.Write(w, binary.LittleEndian, uint16(len(runes)))
	if err != nil {
		return err
	}
	err = binary.Write(w, binary.LittleEndian, runes)
	return nil
}

func writeChar(w io.Writer, char Char) error {
	// vertex amount
	err := binary.Write(w, binary.LittleEndian, uint16(len(char.Vertices)))
	if err != nil {
		return err
	}
	// vertices
	err = binary.Write(w, binary.LittleEndian, char.Vertices)
	if err != nil {
		return err
	}
	// index amount
	err = binary.Write(w, binary.LittleEndian, uint32(len(char.Indices)))
	if err != nil {
		return err
	}
	// indices
	err = binary.Write(w, binary.LittleEndian, char.Indices)
	if err != nil {
		return err
	}
	// advance
	err = binary.Write(w, binary.LittleEndian, char.Advance)
	if err != nil {
		return err
	}
	return nil
}

func FromBinary(r io.Reader) (*Font, error) {
	var amt uint16
	err := binary.Read(r, binary.LittleEndian, &amt)
	if err != nil {
		return nil, err
	}
	// gets a list of the runes contained in the font
	runes, err := readRunes(r, amt)
	if err != nil {
		return nil, err
	}
	// gets the meshes of the runes
	chars, err := readChars(r, runes)
	if err != nil {
		return nil, err
	}
	return &Font{chars}, nil
}

func readChars(r io.Reader, runes []rune) (map[rune]Char, error) {
	chars := make(map[rune]Char)
	for _, run := range runes {
		char, err := readChar(r)
		if err != nil {
			return nil, err
		}
		chars[run] = *char
	}
	return chars, nil
}

func readChar(r io.Reader) (*Char, error) {
	// amount of verts
	var vertAmt uint16
	err := binary.Read(r, binary.LittleEndian, &vertAmt)
	if err != nil {
		return nil, err
	}
	// verts
	vertices := make([][2]float32, vertAmt)
	err = binary.Read(r, binary.LittleEndian, &vertices)
	if err != nil {
		return nil, err
	}
	// amount of indices
	var indAmt uint32
	err = binary.Read(r, binary.LittleEndian, &indAmt)
	if err != nil {
		return nil, err
	}
	// indices
	indices := make([]uint16, indAmt)
	err = binary.Read(r, binary.LittleEndian, &indices)
	if err != nil {
		return nil, err
	}
	// advance
	var advance float32
	err = binary.Read(r, binary.LittleEndian, &advance)
	if err != nil {
		return nil, err
	}
	return &Char{vertices, indices, advance}, nil
}

func readRunes(r io.Reader, amt uint16) ([]rune, error) {
	runes := make([]rune, amt)
	err := binary.Read(r, binary.LittleEndian, &runes)
	if err != nil {
		return nil, err
	}
	return runes, nil
}
