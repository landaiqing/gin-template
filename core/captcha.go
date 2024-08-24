package core

import (
	"github.com/golang/freetype/truetype"
	"github.com/wenlng/go-captcha-assets/bindata/chars"
	"github.com/wenlng/go-captcha-assets/resources/fonts/fzshengsksjw"
	"github.com/wenlng/go-captcha-assets/resources/images"
	"github.com/wenlng/go-captcha-assets/resources/shapes"
	"github.com/wenlng/go-captcha-assets/resources/tiles"
	"github.com/wenlng/go-captcha/v2/base/option"
	"github.com/wenlng/go-captcha/v2/click"
	"github.com/wenlng/go-captcha/v2/rotate"
	"github.com/wenlng/go-captcha/v2/slide"
	"schisandra-cloud-album/global"
)

func InitCaptcha() {
	initRotateCaptcha()
}

// initTextCaptcha 初始化点选验证码
func initTextCaptcha() {
	builder := click.NewBuilder(
		click.WithRangeLen(option.RangeVal{Min: 4, Max: 6}),
		click.WithRangeVerifyLen(option.RangeVal{Min: 2, Max: 4}),
		click.WithRangeThumbColors([]string{
			"#1f55c4",
			"#780592",
			"#2f6b00",
			"#910000",
			"#864401",
			"#675901",
			"#016e5c",
		}),
		click.WithRangeColors([]string{
			"#fde98e",
			"#60c1ff",
			"#fcb08e",
			"#fb88ff",
			"#b4fed4",
			"#cbfaa9",
			"#78d6f8",
		}),
	)

	// fonts
	fonts, err := fzshengsksjw.GetFont()
	if err != nil {
		global.LOG.Fatalln(err)
	}

	// background images
	imgs, err := images.GetImages()
	if err != nil {
		global.LOG.Fatalln(err)
	}

	// thumb images
	//thumbImages, err := thumbs.GetThumbs()
	//if err != nil {
	//	log.Fatalln(err)
	//}

	// set resources
	builder.SetResources(
		click.WithChars(chars.GetChineseChars()),
		//click.WithChars([]string{
		//	"1A",
		//	"5E",
		//	"3d",
		//	"0p",
		//	"78",
		//	"DL",
		//	"CB",
		//	"9M",
		//}),
		//click.WithChars(chars.GetAlphaChars()),
		click.WithFonts([]*truetype.Font{fonts}),
		click.WithBackgrounds(imgs),
		//click.WithThumbBackgrounds(thumbImages),
	)
	global.TextCaptcha = builder.Make()

	// ============================

	builder.Clear()
	builder.SetOptions(
		click.WithRangeLen(option.RangeVal{Min: 4, Max: 6}),
		click.WithRangeVerifyLen(option.RangeVal{Min: 2, Max: 4}),
		click.WithRangeThumbColors([]string{
			"#4a85fb",
			"#d93ffb",
			"#56be01",
			"#ee2b2b",
			"#cd6904",
			"#b49b03",
			"#01ad90",
		}),
	)
	builder.SetResources(
		click.WithChars(chars.GetChineseChars()),
		click.WithFonts([]*truetype.Font{fonts}),
		click.WithBackgrounds(imgs),
	)
	global.LightTextCaptcha = builder.Make()
}

// initClickShapeCaptcha 初始化点击形状验证码
func initClickShapeCaptcha() {
	builder := click.NewBuilder(
		click.WithRangeLen(option.RangeVal{Min: 3, Max: 6}),
		click.WithRangeVerifyLen(option.RangeVal{Min: 2, Max: 3}),
		click.WithRangeThumbBgDistort(1),
		click.WithIsThumbNonDeformAbility(true),
	)

	// shape
	// click.WithUseShapeOriginalColor(false) -> Random rewriting of graphic colors
	shapeMaps, err := shapes.GetShapes()
	if err != nil {
		global.LOG.Fatalln(err)
	}

	// background images
	imgs, err := images.GetImages()
	if err != nil {
		global.LOG.Fatalln(err)
	}

	// set resources
	builder.SetResources(
		click.WithShapes(shapeMaps),
		click.WithBackgrounds(imgs),
	)
	global.ClickShapeCaptcha = builder.MakeWithShape()
}

// initSlideCaptcha 初始化滑动验证码
func initsSlideCaptcha() {
	builder := slide.NewBuilder(
		//slide.WithGenGraphNumber(2),
		slide.WithEnableGraphVerticalRandom(true),
	)

	// background images
	imgs, err := images.GetImages()
	if err != nil {
		global.LOG.Fatalln(err)
	}

	graphs, err := tiles.GetTiles()
	if err != nil {
		global.LOG.Fatalln(err)
	}

	var newGraphs = make([]*slide.GraphImage, 0, len(graphs))
	for i := 0; i < len(graphs); i++ {
		graph := graphs[i]
		newGraphs = append(newGraphs, &slide.GraphImage{
			OverlayImage: graph.OverlayImage,
			MaskImage:    graph.MaskImage,
			ShadowImage:  graph.ShadowImage,
		})
	}

	// set resources
	builder.SetResources(
		slide.WithGraphImages(newGraphs),
		slide.WithBackgrounds(imgs),
	)

	global.SlideCaptcha = builder.Make()
}

// initRotateCaptcha 初始化旋转验证码
func initRotateCaptcha() {
	builder := rotate.NewBuilder(rotate.WithRangeAnglePos([]option.RangeVal{
		{Min: 20, Max: 330},
	}))

	// background images
	imgs, err := images.GetImages()
	if err != nil {
		global.LOG.Fatalln(err)
	}

	// set resources
	builder.SetResources(
		rotate.WithImages(imgs),
	)

	global.RotateCaptcha = builder.Make()
}

// initSlideRegionCaptcha 初始化滑动区域验证码
func initSlideRegionCaptcha() {
	builder := slide.NewBuilder(
		slide.WithGenGraphNumber(2),
		slide.WithEnableGraphVerticalRandom(true),
	)

	// background image
	imgs, err := images.GetImages()
	if err != nil {
		global.LOG.Fatalln(err)
	}

	graphs, err := tiles.GetTiles()
	if err != nil {
		global.LOG.Fatalln(err)
	}
	var newGraphs = make([]*slide.GraphImage, 0, len(graphs))
	for i := 0; i < len(graphs); i++ {
		graph := graphs[i]
		newGraphs = append(newGraphs, &slide.GraphImage{
			OverlayImage: graph.OverlayImage,
			MaskImage:    graph.MaskImage,
			ShadowImage:  graph.ShadowImage,
		})
	}

	// set resources
	builder.SetResources(
		slide.WithGraphImages(newGraphs),
		slide.WithBackgrounds(imgs),
	)

	global.SlideRegionCaptcha = builder.MakeWithRegion()
}
