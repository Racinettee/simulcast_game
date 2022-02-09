package component

import ebi "github.com/hajimehoshi/ebiten/v2"

type Renderer interface {
	RenderItem(img *ebi.Image, options *ebi.DrawImageOptions)
	GetTranslation(x, y float64) *ebi.DrawImageOptions
}

type Renderable interface {
	Render(Renderer)
}
