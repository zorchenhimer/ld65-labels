package main

// FCEUX helper shit

import (
    "fmt"
    "os"
    "strings"
)

type FCEUX struct {
    name        string
    ram         *os.File
    pages       map[int]*os.File
}

func NewFceux(basename string) (*FCEUX, error) {
    r, err := os.Create(strings.Replace(basename, ".db", ".ram.nl", -1))
    if err != nil {
        return nil, err
    }

    return &FCEUX{
        name:   basename,
        ram:    r,
        pages:  make(map[int]*os.File),
    }, nil
}

func (f *FCEUX) Close() {
    f.ram.Close()
    for _, page := range f.pages {
        page.Close()
    }
}

func (f *FCEUX) Write(sym *Symbol, seg *Segment) error {
    var file *os.File
    var err error

    if seg.IsRam() {
        file = f.ram
    } else if file, err = f.getFile(seg.PageID()); err != nil {
        return err
    }

    fmt.Fprintf(file, "$%04X#%s%s#\n", sym.Value, sym.Parent, sym.Name)
    return nil
}

func (f *FCEUX) getFile(pageId int) (*os.File, error) {
    if file, ok := f.pages[pageId]; ok {
        return file, nil
    }
    newFile, err := os.Create(
        strings.Replace(
            f.name,
            ".db",
            fmt.Sprintf(".%d.nl", pageId),
            -1,
        ))
    if err != nil {
        return nil, err
    }
    f.pages[pageId] = newFile
    return newFile, nil
}
