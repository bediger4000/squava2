package mover

import (
	"bytes"
	"fmt"
	"log"
	"os"
)

type Mvr struct {
	moves       [][]byte
	moveCounter int
	count       int
}

func NewFromFile(fileName string) *Mvr {
	buf, err := os.ReadFile(fileName)
	if err != nil {
		log.Print(err)
		return nil
	}
	return NewFromBuffer(buf)
}

func NewFromBuffer(buffer []byte) *Mvr {
	var moves [][]byte
	lines := bytes.Count(buffer, []byte{'\n'})
	if lines > 1 {
		// assume 1 move per line
		moves = bytes.Split(bytes.TrimSpace(buffer), []byte{'\n'})
	} else {
		// assume all moves on one line
		moves = bytes.Split(bytes.TrimSpace(buffer), []byte{' '})
	}
	return &Mvr{
		moves:       moves,
		moveCounter: 0,
		count:       len(moves),
	}
}

func (m *Mvr) Next() (int, int, bool) {

	if m.moveCounter == m.count {
		return 0, 0, false
	}

	move := m.moves[m.moveCounter]
	defer func() { m.moveCounter++ }()

	fields := bytes.Split(move, []byte{','})

	if len(fields) != 2 {
		fmt.Fprintf(os.Stderr, "Move %d, %q, problem\n", m.moveCounter, string(move))
		return 0, 0, false
	}

	if len(fields[0]) == 2 {
		fields[0] = fields[0][0:1]
	}
	if len(fields[1]) == 2 {
		fields[1] = fields[1][0:1]
	}

	if fields[0][0] < '0' || fields[0][0] > '4' {
		fmt.Fprintf(os.Stderr, "Move %d, %q, problem with 1st field\n", m.moveCounter, string(move))
		return 0, 0, false
	}

	x := int(fields[0][0] - '0')

	if fields[1][0] < '0' || fields[1][0] > '4' {
		fmt.Fprintf(os.Stderr, "Move %d, %q, problem with 2nd field\n", m.moveCounter, string(move))
		return 0, 0, false
	}

	y := int(fields[1][0] - '0')

	return x, y, true
}
