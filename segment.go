package main

import (
    "fmt"
    "strconv"
    "strings"
)

type Segment struct {
    Id      int     // unique segment ID.  Symbols link to this.
    Name    string  // name of the segment ("PAGE0", "BSS", etc).
    Start   int     // start of address space
    Size    int     // size of the assembled segment
    Type    string  //

    OutputFile      string  // File this segment will be linked into
    OutputOffset    int     // Offset in the output file
}

func (s *Segment) String() string {
    return fmt.Sprintf("[%d] %s type:%s start:$%04X Size:$%04X [%s @ $%04X]",
        s.Id, s.Name, s.Type, s.Start, s.Size, s.OutputFile, s.OutputOffset)
}

// IsSegment returns true if the given line describes a Segment
func IsSegment(line string) bool {
    if line[0:4] == "seg\t" {
        return true
    }
    return false
}

// ParseSegment returns a *Segment given a line
func ParseSegment(line string) (*Segment, error) {
    if len(line) < 4 || line[0:4] != "seg\t" {
        return nil, fmt.Errorf("Invalid segment line: %q", line)
    }

    items := strings.Split(line[4:], ",")
    seg := &Segment{}
    for _, item := range items {
        keyval := strings.SplitN(item, "=", 2)
        switch keyval[0] {
            // Symbols reference this ID.
            case "id":
                id, err := strconv.ParseInt(keyval[1], 10, 32)
                if err != nil {
                    return nil, fmt.Errorf("Invalid segment Id: %q: %s", keyval[1], err)
                }
                seg.Id = int(id)

            // Address of the start of the segment
            case "start":
                val, err := strconv.ParseInt(keyval[1], 0, 32)
                if err != nil {
                    return nil, fmt.Errorf("Invalid segment Start value: %q: %s", keyval[1], err)
                }
                seg.Start = int(val)

            // Unneeded, but parsed anyway
            case "name":
                seg.Name = keyval[1][1:len(keyval[1]) - 1]

            // Also unneeded, but parsed anyway.
            case "size":
                val, err := strconv.ParseInt(keyval[1], 0, 32)
                if err != nil {
                    return nil, fmt.Errorf("Invalid segment Size value: %q: %s", keyval[1], err)
                }
                seg.Size = int(val)

            case "type":
                seg.Type = keyval[1]

            // Output file offset.  Used to calculate a symbol's location in the ROM file.
            case "ooffs":
                val, err := strconv.ParseInt(keyval[1], 0, 32)
                if err != nil {
                    return nil, fmt.Errorf("Invalid segment OutputOffset value: %q: %s", keyval[1], err)
                }
                seg.OutputOffset = int(val)

            // Output filename.  If blank, it's not written to the ROM and is
            // ignored later on (eg "ZEROPAGE" and "BSS" segments)
            case "oname":
                seg.OutputFile = keyval[1][1:len(keyval[1]) - 1]
        }
    }

    return seg, nil
}

