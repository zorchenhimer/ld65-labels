SRC = main.go segment.go symbol.go

.PHONY: clean

all: $(SRC)
	go build

clean:
	-rm ld65-labels.exe ld65-labels
