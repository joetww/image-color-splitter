package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"os"
	"strconv"
	"strings"

	"github.com/nfnt/resize"
)

// 計算每個色碼的出現次數
func dominantColor(img image.Image) color.Color {
	colorCount := make(map[color.Color]int)

	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := img.At(x, y)
			colorCount[c]++
		}
	}

	// 尋找最多的顏色
	var maxColor color.Color
	maxCount := 0
	for c, count := range colorCount {
		if count > maxCount {
			maxCount = count
			maxColor = c
		}
	}
	return maxColor
}

// 將圖片等分成多個區塊
func splitImage(img image.Image, rows, cols int) []image.Image {
	bounds := img.Bounds()
	width, height := bounds.Max.X/cols, bounds.Max.Y/rows

	var subImages []image.Image
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			subRect := image.Rect(x*width, y*height, (x+1)*width, (y+1)*height)
			subImg := image.NewRGBA(subRect)
			draw.Draw(subImg, subRect, img, subRect.Min, draw.Src)
			subImages = append(subImages, subImg)
		}
	}
	return subImages
}

// 根據主導顏色填充圖片
func fillImageWithColors(colors []color.Color, imgWidth, imgHeight, rows, cols int) *image.RGBA {
	cellWidth := imgWidth / cols
	cellHeight := imgHeight / rows

	newImage := image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))

	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			rect := image.Rect(x*cellWidth, y*cellHeight, (x+1)*cellWidth, (y+1)*cellHeight)
			draw.Draw(newImage, rect, &image.Uniform{colors[y*cols+x]}, image.Point{}, draw.Src)
		}
	}

	return newImage
}

func main() {
	// 解析命令行參數
	inputPath := flag.String("input", "image.jpg", "Path to the input image file.")
	outputPath := flag.String("output", "output.png", "Path to the output image file.")
	jsonPath := flag.String("json", "colors.json", "Path to the JSON file for color information.")
	grid := flag.String("grid", "4x8", "Grid size in rowsxcols format.")
	flag.Parse()

	// 解析 grid 參數
	gridParts := strings.Split(*grid, "x")
	if len(gridParts) != 2 {
		fmt.Println("Invalid grid format. Use rowsxcols format.")
		return
	}
	rows, err := strconv.Atoi(gridParts[0])
	if err != nil {
		fmt.Println("Invalid rows value:", err)
		return
	}
	cols, err := strconv.Atoi(gridParts[1])
	if err != nil {
		fmt.Println("Invalid cols value:", err)
		return
	}

	// 開啟圖片
	file, err := os.Open(*inputPath)
	if err != nil {
		fmt.Println("Error opening image:", err)
		return
	}
	defer file.Close()

	// 解碼圖片
	img, err := jpeg.Decode(file)
	if err != nil {
		fmt.Println("Error decoding image:", err)
		return
	}

	// 調整圖片大小，確保圖片的維度可以均分
	img = resize.Resize(320, 320, img, resize.Lanczos3)

	// 將圖片等分成多個區塊
	subImages := splitImage(img, rows, cols)

	// 對每一格計算最多的顏色
	var dominantColors []color.Color
	colorInfo := make([]map[string]interface{}, 0)
	for i, subImg := range subImages {
		dominant := dominantColor(subImg)
		dominantColors = append(dominantColors, dominant)
		r, g, b, _ := dominant.RGBA()
		colorHex := fmt.Sprintf("#%02x%02x%02x", r>>8, g>>8, b>>8)
		colorInfo = append(colorInfo, map[string]interface{}{
			"subimage": i + 1,
			"color":    colorHex,
		})
		fmt.Printf("Sub-image %d dominant color: %s\n", i+1, colorHex)
	}

	// 輸出顏色資訊到 JSON 文件
	if *jsonPath != "" {
		jsonFile, err := os.Create(*jsonPath)
		if err != nil {
			fmt.Println("Error creating JSON file:", err)
			return
		}
		defer jsonFile.Close()

		encoder := json.NewEncoder(jsonFile)
		encoder.SetIndent("", "  ")
		err = encoder.Encode(colorInfo)
		if err != nil {
			fmt.Println("Error encoding JSON file:", err)
			return
		}
	}

	// 根據主導顏色填充新的圖片
	newImg := fillImageWithColors(dominantColors, 320, 320, rows, cols)

	// 保存新的圖片
	outputFile, err := os.Create(*outputPath)
	if err != nil {
		fmt.Println("Error creating output file:", err)
		return
	}
	defer outputFile.Close()

	err = png.Encode(outputFile, newImg)
	if err != nil {
		fmt.Println("Error encoding output image:", err)
		return
	}

	fmt.Println("New image with dominant colors created:", *outputPath)
}

