package pkg

import (
	"bytes"
	"errors"
	"fmt"
)

func isOptional(field *[]byte) bool {
	if field == nil || len(*field) == 0 {
		return false
	}

	if (*field)[len(*field)-1] != '?' {
		return false
	}

	(*field) = (*field)[:len(*field)-1]
	return true
}

func ParseLine(line []byte) ([]*Rule, error) {
	if len(line) == 0 {
		return nil, errors.New("empty line")
	}

	fields := bytes.Fields(line)

	if len(fields) < 3 {
		return nil, fmt.Errorf("expected at least 3 fields, got %d", len(fields))
	}

	if !bytes.Equal(fields[1], []byte(":")) {
		return nil, fmt.Errorf("expected colon, got %q", fields[1])
	}

	var indices []int

	for i, field := range fields[2:] {
		if bytes.Equal(field, []byte("|")) {
			indices = append(indices, i+2)
		}
	}

	// lhs := fields[0]

	if len(indices) == 0 {
		/* fields = fields[2:]

		var rhss [][]byte

		for i := 0; i < len(fields); i++ {
			ok := isOptional(&fields[i])
			if !ok {
				continue
			}

			new_fields := make([][]byte, 0, len(fields))
			for i := 0; i < len(fields); i++ {
				new_field := make([]byte, len(fields[i]))
				copy(new_field, fields[i])
				new_fields = append(new_fields, new_field)
			}
		}

		for i, field := range fields[2:] {
			isOptional(&field)
		} */

		r, err := NewRule(fields[0], fields[2:])
		return []*Rule{r}, err
	}

	left := 2

	var rules []*Rule

	for _, idx := range indices {
		r, err := NewRule(fields[0], fields[left:idx])
		if err != nil {
			return nil, err
		}

		rules = append(rules, r)
		left = idx + 1
	}

	r, err := NewRule(fields[0], fields[left:])
	if err != nil {
		return nil, err
	}

	return append(rules, r), nil
}

func Parse(data []byte) ([]*Rule, error) {
	// results, err := slgp.Parse(data)

	lines := bytes.Split(data, []byte{';'})

	var rules []*Rule

	for i, line := range lines {
		line = bytes.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		rs, err := ParseLine(line)
		if err != nil {
			return nil, fmt.Errorf("failed to parse line %d: %w", i, err)
		} else if len(rs) > 0 {
			rules = append(rules, rs...)
		}
	}

	return rules, nil
}
