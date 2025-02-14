package ascii

import (
	"errors"
	"fy-novel/internal/model"
	"github.com/nfnt/resize"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"strings"
)

type Yukkuri struct {
	Img    image.Image
	File   *os.File
	Greyer GreyFunction
	Params model.YukkuriParams
}

func NewYukkuri(params model.YukkuriParams) *Yukkuri {
	ykr := &Yukkuri{
		Params: params,
		Greyer: NewGreyHandler(),
	}
	return ykr
}

func (ykr *Yukkuri) TransImgToGrey() (*image.Gray, error) {
	if f, err := os.OpenFile(ykr.Params.ImgPath, os.O_RDONLY, os.ModeTemporary); err != nil {
		return nil, err
	} else {
		defer f.Close()

		ykr.File = f
		tmp := strings.Split(f.Name(), ".")
		suffix := strings.ToLower(tmp[len(tmp)-1])
		var decodeErr error
		switch suffix {
		case "jpg":
			ykr.Img, decodeErr = jpeg.Decode(f)
		case "jpeg":
			ykr.Img, decodeErr = jpeg.Decode(f)
		case "png":
			ykr.Img, decodeErr = png.Decode(f)
		default:
			return nil, errors.New("unsupported image type")
		}
		if decodeErr != nil {
			return nil, decodeErr
		}

		rtg := ykr.Img.Bounds()
		log.Printf("image size %d*%d", ykr.Img.Bounds().Max.X, ykr.Img.Bounds().Max.Y)

		// resize

		if ykr.Params.AscWidth > 0 && ykr.Params.AscHeight > 0 {
			ykr.Img = resize.Resize(uint(ykr.Params.AscWidth), uint(ykr.Params.AscHeight), ykr.Img, resize.Lanczos3)
			log.Printf("resize image to %d*%d\n", ykr.Params.AscWidth, ykr.Params.AscHeight)
		}

		// gray handle
		grey := image.NewGray(ykr.Img.Bounds())
		for i := rtg.Min.X; i < rtg.Max.X; i++ {
			for j := rtg.Min.Y; j < rtg.Max.Y; j++ {
				grey.SetGray(ykr.Greyer.GreyFunc(i, j, ykr.Img, ykr.Params.Threshold))
			}
		}

		// write tmp gray img
		if ykr.Params.TmpImgName != "" {
			w := NewImgWriter()
			log.Printf("output tmp image file [%s]", ykr.Params.TmpImgName)
			if err := w.writeImg(ykr.Params.TmpImgName, grey); err != nil {
				log.Panic(err)
			}
		}

		return grey, nil
	}
}

func (ykr *Yukkuri) TransImgToAsc(converter Asc11Converter) {

	converter.Convert()

	w := NewImgWriter()
	if err := w.writeAscii(ykr.Params.Filename, converter.CharMap()); err != nil {
		log.Panic(err)
	}
	log.Printf("output is complete [%s]", ykr.Params.Filename)
}
