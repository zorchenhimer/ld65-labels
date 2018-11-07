package main

import (
    "fmt"
    "strconv"
    "strings"
)

type Symbol struct {
    Id          int     // unique label ID
    Name        string  // Label
    AddrSize    string  // absolute, zeropage, etc
    Parent      int     // ?
    Defined     []int   // "def"
    References  []int   // "ref"
    Value       int     // Address of label
    Segment     int     // ID of segment
}

func (s *Symbol) String() string {
    d := []string{}
    if len(s.Defined) > 0 {
        for _, i := range s.Defined {
            d = append(d, fmt.Sprintf("%d", i))
        }
    }

    r := []string{}
    if len(s.References) > 0 {
        for _, i := range s.References {
            r = append(r, fmt.Sprintf("%d", i))
        }
    }

    return fmt.Sprintf("[%d] %s AddrSize:%s Defined:%s References:%s Value:$%X Segment:%d",
        s.Id, s.Name, s.AddrSize, d, r, s.Value, s.Segment)
}

// IsSymbol returns true if the given line describes a symbol (Label)
func IsSymbol(line string) bool {
    if line[0:4] == "sym\t" {
        return true
    }
    return false
}

// ParseSymbol returns a *Symbol of the given line
func ParseSymbol(line string) (*Symbol, error) {
    if len(line) < 4 || line[0:4] != "sym\t" {
        return nil, fmt.Errorf("Invalid symbol line: %q", line)
    }

    items := strings.Split(line[4:], ",")
    sym := &Symbol{Segment: -1}
    for _, item := range items {
        keyval := strings.SplitN(item, "=", 2)
        switch keyval[0] {
            // Not really needed, but parsed anyway.
            case "id":
                val, err := strconv.ParseInt(keyval[1], 10, 32)
                if err != nil {
                    return nil, fmt.Errorf("Invalid symbol Id: %q: %s", keyval[1], err)
                }
                sym.Id = int(val)

            // The label name that we need
            case "name":
                sym.Name = keyval[1][1:len(keyval[1]) - 1]

            // Not really needed, but parsed anyway.
            case "addrsize":
                sym.AddrSize = keyval[1]

            // These are both lists of numbers delemited by a plus sign.  Aren't
            // needed for my purposes, but I parse it anyway.
            case "ref", "def":
                vals := strings.Split(keyval[1], "+")
                if len(vals) > 0 {
                    lst := []int{}

                    for _, v := range vals {
                        val, err := strconv.ParseInt(v, 10, 32)
                        if err != nil {
                            return nil, fmt.Errorf("Invalid symbol Reference value: %q: %s", v, err)
                        }
                        lst = append(lst, int(val))
                    }

                    if keyval[0] == "ref" {
                        sym.References = lst
                    } else {
                        sym.Defined = lst
                    }
                }

            // Address of the label
            case "val":
                val, err := strconv.ParseInt(keyval[1], 0, 32)
                if err != nil {
                    return nil, fmt.Errorf("Invalid symbol Value value: %q: %s", keyval[1], err)
                }
                sym.Value = int(val)

            // Segment that the label belongs to
            case "seg":
                val, err := strconv.ParseInt(keyval[1], 10, 32)
                if err != nil {
                    return nil, fmt.Errorf("Invalid symbol Segment value: %q: %s", keyval[1], err)
                }
                sym.Segment = int(val)

        }
    }

    return sym, nil
}
