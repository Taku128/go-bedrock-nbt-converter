package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ntaku256/go-bedrock-nbt-converter/pkg/mcstruct"
	"github.com/ntaku256/go-bedrock-nbt-converter/pkg/mcworld"
)

func main() {
	var outPath string
	var minX, maxX, minY, maxY, minZ, maxZ int
	var verbose bool

	flag.StringVar(&outPath, "o", "", "Output file path (default: <input>.nbt)")
	flag.StringVar(&outPath, "out", "", "Output file path (default: <input>.nbt)")

	flag.IntVar(&minX, "x", math.MinInt32, "Minimum X coordinate (only for .mcworld)")
	flag.IntVar(&minX, "min-x", math.MinInt32, "Minimum X coordinate (only for .mcworld)")
	flag.IntVar(&maxX, "X", math.MaxInt32, "Maximum X coordinate (only for .mcworld)")
	flag.IntVar(&maxX, "max-x", math.MaxInt32, "Maximum X coordinate (only for .mcworld)")

	flag.IntVar(&minY, "y", -64, "Minimum Y coordinate (only for .mcworld)")
	flag.IntVar(&minY, "min-y", -64, "Minimum Y coordinate (only for .mcworld)")
	flag.IntVar(&maxY, "Y", 320, "Maximum Y coordinate (only for .mcworld)")
	flag.IntVar(&maxY, "max-y", 320, "Maximum Y coordinate (only for .mcworld)")

	flag.IntVar(&minZ, "z", math.MinInt32, "Minimum Z coordinate (only for .mcworld)")
	flag.IntVar(&minZ, "min-z", math.MinInt32, "Minimum Z coordinate (only for .mcworld)")
	flag.IntVar(&maxZ, "Z", math.MaxInt32, "Maximum Z coordinate (only for .mcworld)")
	flag.IntVar(&maxZ, "max-z", math.MaxInt32, "Maximum Z coordinate (only for .mcworld)")

	flag.BoolVar(&verbose, "v", false, "Enable verbose output")
	flag.BoolVar(&verbose, "verbose", false, "Enable verbose output")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: bedrock-nbt-converter <inputPath> [options]\n\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		flag.Usage()
		os.Exit(1)
	}

	inputPath := args[0]
	ext := strings.ToLower(filepath.Ext(inputPath))

	if outPath == "" {
		outPath = strings.TrimSuffix(inputPath, ext) + ".nbt"
	}

	if verbose {
		fmt.Printf("Input: %s\n", inputPath)
		fmt.Printf("Output: %s\n", outPath)
		if ext == ".mcworld" {
			fmt.Printf("Bounding Box: X[%d, %d] Y[%d, %d] Z[%d, %d]\n", minX, maxX, minY, maxY, minZ, maxZ)
		}
	}

	startTime := time.Now()
	var nbtData []byte
	var err error

	if ext == ".mcworld" {
		opts := &mcworld.ConvertOptions{
			MinX: int32(minX), MaxX: int32(maxX),
			MinY: int32(minY), MaxY: int32(maxY),
			MinZ: int32(minZ), MaxZ: int32(maxZ),
			Dimension: 0,
		}
		nbtData, err = mcworld.ConvertMcworld(inputPath, opts)
	} else if ext == ".mcstructure" {
		nbtData, err = mcstruct.ConvertMcstructure(inputPath)
	} else {
		fmt.Fprintf(os.Stderr, "Error: Unknown file extension '%s'. Must be .mcworld or .mcstructure.\n", ext)
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Conversion failed: %v\n", err)
		os.Exit(1)
	}

	err = os.WriteFile(outPath, nbtData, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write output file: %v\n", err)
		os.Exit(1)
	}

	elapsed := time.Since(startTime)
	if verbose {
		fmt.Printf("Successfully saved NBT to %s in %s\n", outPath, elapsed)
	} else {
		fmt.Printf("Done in %s. Saved to %s.\n", elapsed, outPath)
	}
}
