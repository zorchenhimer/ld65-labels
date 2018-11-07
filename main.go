package main

// TODO: Add option to output FCEUX symbol files

import (
    "io/ioutil"
    "fmt"
    "os"
    "strings"
)

var (
    verbose     bool
    inputName   string
    outputMesen bool
    outputFceux bool
)

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Missing input file")
        os.Exit(1)
    }

    inputname := os.Args[1]

    file, err := ioutil.ReadFile(inputname)
    if err != nil {
        fmt.Println("Unable to read input file: ", err)
        os.Exit(1)
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
        os.Exit(1)
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
            abs = sym.Value - seg.Start + seg.OutputOffset - 16
        }

        if abs < 0 {
            fmt.Printf("ERROR: absolute address is negative!\n  Label:%q Segment:%q abs:%d sym.Value:%d seg.Start:%d seg.OutputOffset:%d\n",
                sym.Name,
                seg.Name,
                abs,
                sym.Value,
                seg.Start,
                seg.OutputOffset,
            )
            os.Exit(1)
        }

        // Type (PRG/RAM), Address (loaded into CPU space), Name
        fmt.Fprintf(outfile, "%s:%04X:%s\n", t, abs, sym.Name)
        count += 1
    }

    // Report how many we found
    fmt.Printf("  Exported labels: %d\n", count)
}
