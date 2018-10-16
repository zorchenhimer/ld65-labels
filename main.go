package main

// TODO: Add option to output FCEUX symbol files

import (
    "io/ioutil"
    "fmt"
    "os"
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

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Missing input file")
        return
    }

    inputname := os.Args[1]

    file, err := ioutil.ReadFile(inputname)
    if err != nil {
        fmt.Println("Unable to read input file: ", err)
        return
    }

    lines := strings.Split(strings.Replace(string(file), "\r", "",  -1), "\n")
    segments := []*Segment{}
    symbols := []*Symbol{}

    // Read input file and parse all segments and symbols
    for _, line := range lines {
        // Ignore invalid lines
        if len(line) < 4 {
            continue
        }

        if IsSegment(line) {
            seg, err := ParseSegment(line)
            if err != nil {
                fmt.Println(err)
            }
            segments = append(segments, seg)

        } else if IsSymbol(line) {
            sym, err := ParseSymbol(line)
            if err != nil {
                fmt.Println(err)
            }
            symbols = append(symbols, sym)
        }
    }

    // Create the .mlb output file for Mesen
    // TODO: output FCEUX's format too
    outfile, err := os.Create(strings.Replace(inputname, ".nes.db", ".mlb", -1))
    if err != nil {
        fmt.Println(err)
        return
    }
    defer outfile.Close()

    count := 0  // Just used to show the user at the end
    for _, sym := range symbols {
        // Ignore symbols that don't have a segment (usually internal symbols)
        if sym.Segment == -1 {
            continue
        }

        // Find the segment for the current symbol
        var seg *Segment
        for _, s := range segments {
            if s.Id == sym.Segment {
                seg = s
                break
            }
        }

        // skip labels that don't have a segment
        if seg == nil {
            continue
        }

        abs := sym.Value

        t := "P"
        // Variable addresses don't need to be adjusted.
        if len(seg.OutputFile) == 0 {
            t = "R"
        } else {
            // minus start of PRG, minus iNES header, plus output file offset
            abs -= seg.Start + seg.OutputOffset - 16
        }

        // Type (PRG/RAM), Address (loaded into CPU space), Name
        fmt.Fprintf(outfile, "%s:%04X:%s\n", t, abs, sym.Name)
        count += 1
    }

    // Report how many we found
    fmt.Printf("  Exported labels: %d\n", count)
}
