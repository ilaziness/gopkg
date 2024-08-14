package main

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"os"

	"github.com/ilaziness/gopkg/image/resize"
	qrcode "github.com/skip2/go-qrcode"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/font/gofont/goitalic"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

func s() []int {
	log.Println(123)
	return []int{5, 6, 1}
}

func main() {
	// 蒙版图像
	mask()

	///////////////////////////////////////////////////
	// 裁剪，合并并应用蒙版
	mi1, _ := os.Open("testdata/test.png")
	clipi, _ := os.Open("testdata/blue.png")
	mi11, _, _ := image.Decode(mi1)
	c11, _, _ := image.Decode(clipi)

	// 图像文件转成RGBA
	dst1 := image.NewRGBA(mi11.Bounds())
	draw.Draw(dst1, mi11.Bounds(), mi11, image.Pt(0, 0), draw.Src)
	dst2 := image.NewRGBA(mi11.Bounds())
	draw.Draw(dst2, mi11.Bounds(), mi11, image.Pt(0, 0), draw.Src)

	// c11图片裁剪一部分作为src
	src1 := image.NewRGBA(image.Rect(0, 0, 50, 50))
	ret := image.Rect(0, 0, 50, 50)
	draw.Draw(src1, src1.Bounds(), c11, ret.Min, draw.Src)
	// dst：目标图像，通常是一个实现了 draw.Image 接口的对象。
	// r：目标图像上的矩形区域，表示源图像和蒙版图像将被绘制到的区域。
	// src：源图像，即要绘制到目标图像上的图像。
	// sp：源图像的起始点，表示从源图像的哪个点开始绘制。
	// mask：蒙版图像，用于控制源图像的透明度。如果不需要蒙版，可以传入 nil。
	// mp：蒙版图像的起始点，表示从蒙版图像的哪个点开始应用蒙版。
	// op：绘制操作，可以是 draw.Over 或 draw.Src。draw.Over 表示源图像与目标图像混合，draw.Src 表示源图像直接替换目标图像的像素。
	draw.DrawMask(
		dst1,
		dst1.Bounds().Intersect(image.Rect(50, 100, 200, 200)),
		src1,
		image.Pt(0, 0), // 相对于r
		createGradientMask2(100, 100),
		image.Pt(0, 0), // 相对于r
		draw.Over,
	)
	// 保存图像到文件
	file, _ := os.Create("dst1.png")
	defer file.Close()
	png.Encode(file, dst1)
	// 圆形蒙版
	draw.DrawMask(
		dst2,
		dst2.Bounds().Intersect(image.Rect(50, 100, 200, 200)),
		src1,
		image.Pt(0, 0), // 相对于r
		&circle{image.Pt(25, 25), 20},
		image.Pt(0, 0), // 相对于r
		draw.Over,
	)
	file, _ = os.Create("dst2.png")
	defer file.Close()
	png.Encode(file, dst2)

	///////////////////////////////////////////////////
	imgFile, _ := os.OpenFile("imgfile.png", os.O_CREATE|os.O_RDWR, 0755)
	img := image.NewRGBA(image.Rect(0, 0, 300, 300))
	// 填充蓝色
	blue := color.RGBA{0, 0, 255, 255}
	// Uniform 表示一个无限大小统一颜色的图形
	draw.Draw(img, img.Bounds(), &image.Uniform{blue}, image.Pt(0, 0), draw.Src)
	_ = png.Encode(imgFile, img)

	///////////////////////////////////////////////////
	//缩放
	smallImg := resize.Resize(img, img.Bounds(), 100, 100)
	simgf, _ := os.Create("simgf.png")
	_ = png.Encode(simgf, smallImg)

	bigImg := resize.Resize(img, img.Bounds(), 600, 600)
	bigf, _ := os.Create("bigf.png")
	_ = png.Encode(bigf, bigImg)

	testF, _ := os.Open("testdata/test.png")
	tf, ext, err := image.Decode(testF)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("extension:", ext)
	// 缩小
	testfs := resize.Resize(tf, tf.Bounds(), 198, 155)
	tfsf, _ := os.Create("tfsf.png")
	png.Encode(tfsf, testfs)
	// 放大
	testfb := resize.Resize(tf, tf.Bounds(), 497, 410)
	tfbf, _ := os.Create("tfbf.png")
	png.Encode(tfbf, testfb)

	// 缩小50%
	tfn := resize.HalveInplace(tf)
	tf5f, _ := os.Create("test50.png")
	png.Encode(tf5f, tfn)

	///////////////////////////////////////////////////
	// 生成二维码
	qrcode.WriteFile("https://example.com", qrcode.Medium, 200, "qr.png")
	qrcode.WriteColorFile("https://example.com", qrcode.Medium, 200, color.Black, color.White, "qr_color.png")

	///////////////////////////////////////////////////
	// 文字渲染
	imgf := image.NewRGBA(image.Rect(0, 0, 200, 200))
	// 填充背景色
	draw.Draw(imgf, imgf.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)
	// 创建一个字体面
	face := basicfont.Face7x13
	// 创建一个绘制器
	d := &font.Drawer{
		Dst:  imgf,
		Src:  image.NewUniform(color.Black),
		Face: face,
		Dot:  fixed.Point26_6{X: fixed.I(10), Y: fixed.I(50)},
	}
	d.DrawString("Hello, Go!")
	d.DrawString("Hello, Go!")

	// 指定字体文字渲染
	f, _ := opentype.Parse(goitalic.TTF)
	face2, _ := opentype.NewFace(f, &opentype.FaceOptions{
		Size:    32,
		DPI:     72,
		Hinting: font.HintingNone,
	})
	d2 := font.Drawer{
		Dst:  imgf,
		Src:  image.Black,
		Face: face2,
		Dot:  fixed.P(10, 80),
	}
	d2.DrawString("test, test")
	d2.Src = image.NewUniform(color.Gray{0x7F})
	d2.DrawString("ly")

	// 保存图像到文件
	fontf, _ := os.Create("font.png")
	defer fontf.Close()
	png.Encode(fontf, imgf)
}

type circle struct {
	p image.Point

	r int
}

func (c *circle) ColorModel() color.Model {

	return color.AlphaModel

}

func (c *circle) Bounds() image.Rectangle {

	return image.Rect(c.p.X-c.r, c.p.Y-c.r, c.p.X+c.r, c.p.Y+c.r)

}

func (c *circle) At(x, y int) color.Color {

	xx, yy, rr := float64(x-c.p.X)+0.5, float64(y-c.p.Y)+0.5, float64(c.r)

	if xx*xx+yy*yy < rr*rr {

		return color.Alpha{255}

	}

	return color.Alpha{0}

}
