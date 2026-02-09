package service

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"

	xdraw "golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

// Processor - логика обработки изображения
type Processor struct{}

// NewProcessor - инициализация NewProcessor
func NewProcessor() *Processor {
	return &Processor{}
}

// Process - обработка изображения
func (p *Processor) Process(r io.Reader, format string) (*bytes.Buffer, error) {
	var img image.Image
	var err error

	// Декодирование изображения в зависимости от формата
	switch format {
	case "gif":
		// Для GIF читаем все кадры, но используем только первый
		g, err := gif.DecodeAll(r)
		if err != nil {
			return nil, err
		}
		if len(g.Image) == 0 {
			return nil, err
		}
		img = g.Image[0]
	default:
		// Для JPG и PNG используем стандартный image.Decode
		img, _, err = image.Decode(r)
		if err != nil {
			return nil, err
		}
	}

	// Изменяем размер изображения
	resized := resizeImage(img, 800)

	// Добавляем водяной знак
	watermarked := addWatermark(resized, "ImageProcessor")

	// Кодируем обработанное изображение обратно в буфер
	buf := new(bytes.Buffer)
	switch format {
	case "png":
		if err := png.Encode(buf, watermarked); err != nil {
			return nil, err
		}
	case "gif":
		// Сохраняем как статический GIF (один кадр)
		if err := gif.Encode(buf, watermarked, nil); err != nil {
			return nil, err
		}
	default:
		// По умолчанию сохраняем в JPEG
		if err := jpeg.Encode(buf, watermarked, &jpeg.Options{Quality: 85}); err != nil {
			return nil, err
		}
	}

	return buf, nil
}

// resizeImage - масштабирование изображения по ширине
func resizeImage(src image.Image, newWidth int) image.Image {
	bounds := src.Bounds()

	// Вычисляем новую высоту
	ratio := float64(newWidth) / float64(bounds.Dx())
	newHeight := int(float64(bounds.Dy()) * ratio)

	// Создаём новый холст под изменённый размер
	dst := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))

	// Ресайз
	xdraw.ApproxBiLinear.Scale(dst, dst.Bounds(), src, bounds, draw.Over, nil)

	return dst
}

// addWatermark - добавление водяного знака
func addWatermark(img image.Image, text string) image.Image {
	// Создаём RGBA копию изображения для рисования поверх
	rgba := image.NewRGBA(img.Bounds())
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{}, draw.Src)

	// Цвет водяного знака (полупрозрачный белый)
	col := color.RGBA{255, 255, 255, 200}

	// Позиция текста(снизу слева)
	point := fixed.Point26_6{
		X: fixed.I(20),
		Y: fixed.I(rgba.Bounds().Dy() - 20),
	}

	// Рисуем текст на изображении
	d := &font.Drawer{
		Dst:  rgba,
		Src:  image.NewUniform(col),
		Face: basicfont.Face7x13,
		Dot:  point,
	}
	d.DrawString(text)

	return rgba
}
