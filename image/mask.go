package main

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
)

func createGradientMask(width, height int) *image.Alpha {
	// 创建一个 Alpha 图像，用于存储渐变蒙版
	mask := image.NewAlpha(image.Rect(0, 0, width, height))

	// 创建渐变蒙版
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// 计算透明度，从 0 到 255, 255是完全透明
			alpha := uint8(255 * (255 - float64(x)/float64(width)))
			mask.SetAlpha(x, y, color.Alpha{A: alpha})
		}
	}

	return mask
}

func createGradientMask2(width, height int) *image.RGBA {
	// 创建一个 Alpha 图像，用于存储渐变蒙版
	mask := image.NewRGBA(image.Rect(0, 0, width, height))

	// 创建渐变蒙版
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// 计算透明度，从 0 到 255, 255是完全透明
			alpha := uint8(255 * (255 - float64(x)/float64(width)))
			mask.SetRGBA(x, y, color.RGBA{0, 0, 0, alpha})
		}
	}

	return mask
}

func mask() {
	// 创建一个 200x200 的图像
	width, height := 200, 200
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// 创建渐变蒙版
	mask := createGradientMask2(width, height)

	// 将渐变蒙版应用到图像上
	// draw.Over表示将蒙版覆盖到图像上
	draw.DrawMask(img, img.Bounds(), image.NewUniform(color.Black), image.Point{}, mask, image.Point{}, draw.Over)

	// 保存图像到文件
	file, err := os.Create("output.png")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	err = png.Encode(file, img)
	if err != nil {
		panic(err)
	}
}
