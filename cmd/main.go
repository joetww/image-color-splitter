package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"os"

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

// 將圖片等分成32個區塊
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
	// 開啟圖片
	file, err := os.Open("image.jpg")
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

	// 將圖片等分成32格 (4行x8列)
	rows, cols := 4, 8
	subImages := splitImage(img, rows, cols)

	// 對每一格計算最多的顏色
	var dominantColors []color.Color
	for i, subImg := range subImages {
		dominant := dominantColor(subImg)
		dominantColors = append(dominantColors, dominant)
		r, g, b, _ := dominant.RGBA()
		fmt.Printf("Sub-image %d dominant color: #%02x%02x%02x\n", i+1, r>>8, g>>8, b>>8)
	}

	// 根據主導顏色填充新的圖片
	newImg := fillImageWithColors(dominantColors, 320, 320, rows, cols)

	// 保存新的圖片
	outputFile, err := os.Create("output.png")
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

	fmt.Println("New image with dominant colors created: output.png")
}

