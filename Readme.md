程式說明
這個 Go 程式會將圖片等分成多個區塊，計算每個區塊的主導顏色，並根據這些顏色生成一個新的圖片。它支持以下功能：

計算主導顏色：對每個區塊中的顏色進行計數，找出出現次數最多的顏色。
分割圖片：將原始圖片等分成指定數量的區塊（例如 4 行 x 8 列）。
填充新圖片：用計算得到的主導顏色填充一個新的圖片，按照原圖片的分割結構排列顏色。
輸出顏色資訊：將每個區塊的顏色資訊輸出為 JSON 格式的文件。
處理命令行參數：支持指定原圖檔案路徑、生成新圖檔路徑、顏色資訊 JSON 文件的路徑和分割格數。
主要函數：
dominantColor(img image.Image) color.Color：計算圖片的主導顏色。
splitImage(img image.Image, rows, cols int) []image.Image：將圖片等分成指定數量的區塊。
fillImageWithColors(colors []color.Color, imgWidth, imgHeight, rows, cols int) *image.RGBA：用顏色填充新圖片。
main()：程序入口，處理命令行參數、執行圖片處理和保存結果。
編譯方法
確保已安裝 Go：在終端中運行 go version 確認 Go 環境已經安裝。

安裝依賴包： 使用以下命令安裝程序所需的 Go 包：

bash
複製程式碼
go get github.com/nfnt/resize
編譯程式碼： 確保你在專案的根目錄下，即包含 cmd 目錄的目錄。運行以下命令來編譯程式：

bash
複製程式碼
go build -o image-color-splitter cmd/main.go
運行程式： 編譯後，你可以使用以下命令運行程式並指定參數：

bash
複製程式碼
./image-color-splitter -input path/to/your/image.jpg -output path/to/your/output.png -json path/to/your/colors.json -grid 4x8
-input：原圖檔案的路徑。
-output：生成的新圖檔路徑。
-json：顏色資訊 JSON 文件的路徑。
-grid：分割格數（例如 4x8 代表 4 行 8 列）。
