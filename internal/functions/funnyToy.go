package functions

import (
	"context"
	"fy-novel/internal/model"
	"fy-novel/internal/toys/ascii"

	"github.com/sirupsen/logrus"
)

type FunnyToy struct {
	log *logrus.Logger
}

func NewFunnyToy(l *logrus.Logger) *FunnyToy {
	return &FunnyToy{log: l}
}

func (ft *FunnyToy) ConvertToAscii(ctx context.Context, params model.YukkuriParams) error {
	ft.log.Printf("load file [%s]", params.ImgPath)
	ykr := ascii.NewYukkuri(params)
	if gray, err := ykr.TransImgToGrey(); err != nil {
		ft.log.Printf("failed to converted the file into ascii %v\n", err)
	} else {
		converter := ascii.NewAsc11Converter(gray, params.Ykr)
		ykr.TransImgToAsc(converter)
	}
	return nil
}
