package netpbm

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"sort"
)

// PPM is a structure to represent PPM images.
type PPM struct {
	Data        [][]Pixel
	Width, Height int
	MagicNumber string
	Max         uint
}

// Pixel represents a pixel with red (R), green (G), and blue (B) channels.
type Pixel struct {
	R, G, B uint8
}

// Point represents a point in the image.
type Point struct {
	X, Y int
}

// ReadPPM lit une image PPM depuis un fichier et renvoie une structure qui représente l'image.
func ReadPPM(filename string) (*PPM, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// lit le numero magique
	if !scanner.Scan() {
		return nil, fmt.Errorf("failed to read magic number")
	}
	magicNumber := scanner.Text()

	if magicNumber != "P3" && magicNumber != "P6" {
		return nil, fmt.Errorf("unsupported PPM format: %s", magicNumber)
	}

	// passe les commentaires et les lignes vides
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) > 0 && line[0] != '#' {
			break
		}
	}

	// lit la largeur et la hauteur
	if scanner.Err() != nil {
		return nil, fmt.Errorf("error reading dimensions line: %v", scanner.Err())
	}
	dimensions := strings.Fields(scanner.Text())
	if len(dimensions) != 2 {
		return nil, fmt.Errorf("invalid dimensions line")
	}

	width, err := strconv.Atoi(dimensions[0])
	if err != nil {
		return nil, fmt.Errorf("failed to parse width: %v", err)
	}

	height, err := strconv.Atoi(dimensions[1])
	if err != nil {
		return nil, fmt.Errorf("failed to parse height: %v", err)
	}

	// lit la valeur max
	if !scanner.Scan() {
		return nil, fmt.Errorf("failed to read max value")
	}
	maxValue, err := strconv.Atoi(scanner.Text())
	if err != nil {
		return nil, fmt.Errorf("failed to parse max value: %v", err)
	}

	// lit les données
	var data [][]Pixel
	for y := 0; y < height; y++ {
		var row []Pixel
		for x := 0; x < width; x++ {
			var pixel Pixel
			if magicNumber == "P3" {
				// format ASCII
				if _, err := fmt.Fscanf(Scanner, "%d %d %d", &pixel.R, &pixel.G, &pixel.B); err != nil {
					return nil, fmt.Errorf("failed to parse pixel data: %v", err)
				}
			} else {
				// Format binaire (P6)
				var buf [3]byte
				if _, err := file.Read(buf[:]); err != nil {
					return nil, fmt.Errorf("failed to read pixel data: %v", err)
				}
				pixel.R, pixel.G, pixel.B = buf[0], buf[1], buf[2]
			}
			row = append(row, pixel)
		}
		data = append(data, row)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	return &PPM{
		Data:        data,
		Width:       width,
		Height:      height,
		MagicNumber: magicNumber,
		Max:         uint(maxValue),
	}, nil
}

// Size renvoie la largeur et la hauteur de l'image.
func (ppm *PPM) Size() (int, int) {
	return ppm.Width, ppm.Height
}

// At renvoie la valeur du pixel à la position (x, y).
func (ppm *PPM) At(x, y int) Pixel {
	return ppm.Data[y][x]
}

// Set définit la valeur du pixel à la position (x, y).
func (ppm *PPM) Set(x, y int, value Pixel) {
	if x >= 0 && x < ppm.Width && y >= 0 && y < ppm.Height {
		ppm.Data[y][x] = value
	}
}

// Save enregistre l'image PPM dans un fichier et renvoie une erreur s'il y a un problème.
func (ppm *PPM) Save(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	// ecrit le nombre magique
	if _, err := fmt.Fprintf(file, "%s\n", ppm.MagicNumber); err != nil {
		return fmt.Errorf("failed to write magic number: %v", err)
	}

	// ecrit la largeur et la hauteur
	if _, err := fmt.Fprintf(file, "%d %d\n", ppm.Width, ppm.Height); err != nil {
		return fmt.Errorf("failed to write dimensions: %v", err)
	}

	// ecrit la valeur max
	if _, err := fmt.Fprintf(file, "%d\n", ppm.Max); err != nil {
		return fmt.Errorf("failed to write max value: %v", err)
	}

	// ecrit les données
	for y := 0; y < ppm.Height; y++ {
		for x := 0; x < ppm.Width; x++ {
			if ppm.MagicNumber == "P3" {
				// format ASCII
				if _, err := fmt.Fprintf(file, "%d %d %d ", ppm.Data[y][x].R, ppm.Data[y][x].G, ppm.Data[y][x].B); err != nil {
					return fmt.Errorf("failed to write pixel data: %v", err)
				}
			} else {
				// format binaire (P6)
				if _, err := file.Write([]byte{ppm.Data[y][x].R, ppm.Data[y][x].G, ppm.Data[y][x].B}); err != nil {
					return fmt.Errorf("failed to write pixel data: %v", err)
				}
			}
		}
		if ppm.MagicNumber == "P3" {
			// ajoute une nouvelle ligne après chaque ligne au format ASCII
			if _, err := fmt.Fprint(file, "\n"); err != nil {
				return fmt.Errorf("failed to write newline: %v", err)
			}
		}
	}

	return nil
}

// Invert inverse les couleurs de l'image PPM.
func (ppm *PPM) Invert() {
	for y := 0; y < ppm.Height; y++ {
		for x := 0; x < ppm.Width; x++ {
			ppm.Data[y][x].R = ppm.Max - ppm.Data[y][x].R
			ppm.Data[y][x].G = ppm.Max - ppm.Data[y][x].G
			ppm.Data[y][x].B = ppm

		}
	}
}
// Flip retourne l'image PPM horizontalement.
func (ppm *PPM) Flip() {
	for y := 0; y < ppm.Height; y++ {
		for x := 0; x < ppm.Width/2; x++ {
			ppm.Data[y][x], ppm.Data[y][ppm.Width-x-1] = ppm.Data[y][ppm.Width-x-1], ppm.Data[y][x]
		}
	}
}

// Flop fait basculer l'image PPM verticalement.
func (ppm *PPM) Flop() {
	for y := 0; y < ppm.Height/2; y++ {
		ppm.Data[y], ppm.Data[ppm.Height-y-1] = ppm.Data[ppm.Height-y-1], ppm.Data[y]
	}
}

// SetMagicNumber définit le nombre magique de l'image PPM.
func (ppm *PPM) SetMagicNumber(magicNumber string) {
	ppm.MagicNumber = magicNumber
}

// SetMaxValue définit la valeur maximale de l'image PPM.
func (ppm *PPM) SetMaxValue(maxValue uint8) {
	ppm.Max = uint(maxValue)
}

// Rotate90CW fait pivoter l’image PPM de 90° dans le sens des aiguilles d’une montre.
func (ppm *PPM) Rotate90CW() {
	// Créez une nouvelle image RGBA et dessinez dessus les pixels pivotés
	newImage := image.NewRGBA(image.Rect(0, 0, ppm.Height, ppm.Width))

	for y := 0; y < ppm.Height; y++ {
		for x := 0; x < ppm.Width; x++ {
			newImage.Set(y, ppm.Width-x-1, color.RGBA{ppm.Data[y][x].R, ppm.Data[y][x].G, ppm.Data[y][x].B, 255})
		}
	}

	// Mettre à jour les propriétés PPM avec les données pivotées
	ppm.Width, ppm.Height = ppm.Height, ppm.Width
	ppm.Data = make([][]Pixel, ppm.Height)
	for y := 0; y < ppm.Height; y++ {
		ppm.Data[y] = make([]Pixel, ppm.Width)
		for x := 0; x < ppm.Width; x++ {
			r, g, b, _ := newImage.At(y, x).RGBA()
			ppm.Data[y][x] = Pixel{uint8(r), uint8(g), uint8(b)}
		}
	}
}

// ToPGM convertit l'image PPM en PGM.
func (ppm *PPM) ToPGM() *PGM {
	// Créez une nouvelle image PGM avec les mêmes dimensions
	pgm := &PGM{
		Data:        make([][]uint8, ppm.Height),
		Width:       ppm.Width,
		Height:      ppm.Height,
		MagicNumber: "P5",
		Max:         ppm.Max,
	}

	for y := 0; y < ppm.Height; y++ {
		pgm.Data[y] = make([]uint8, ppm.Width)
		for x := 0; x < ppm.Width; x++ {
			// Convertir RVB en niveaux de gris en utilisant la méthode de luminosité
			gray := uint8(0.299*float64(ppm.Data[y][x].R) + 0.587*float64(ppm.Data[y][x].G) + 0.114*float64(ppm.Data[y][x].B))
			pgm.Data[y][x] = gray
		}
	}

	return pgm
}

// ToPBM convertit l'image PPM en PBM.
func (ppm *PPM) ToPBM() *PBM {
	// Créer une nouvelle image PBM avec les mêmes dimensions
	pbm := &PBM{
		Data:        make([][]uint8, ppm.Height),
		Width:       ppm.Width,
		Height:      ppm.Height,
		MagicNumber: "P1",
	}

	for y := 0; y < ppm.Height; y++ {
		pbm.Data[y] = make([]uint8, ppm.Width)
		for x := 0; x < ppm.Width; x++ {
			// Convertir RVB en binaire en utilisant un seuil
			threshold := uint8((uint(ppm.Data[y][x].R) + uint(ppm.Data[y][x].G) + uint(ppm.Data[y][x].B)) / 3)
			if threshold > ppm.Max/2 {
				pbm.Data[y][x] = 1
			} else {
				pbm.Data[y][x] = 0
			}
		}
	}

	return pbm
}

func (ppm *PPM) DrawLine(p1, p2 Point, color Pixel) {
	dx := p2.X - p1.X
	dy := p2.Y - p1.Y

	if dx == 0 && dy == 0 {
		ppm.Set(p1.X, p1.Y, color)
		return
	}

	xSign, ySign := 1, 1

	if dx < 0 {
		dx = -dx
		xSign = -1
	}

	if dy < 0 {
		dy = -dy
		ySign = -1
	}

	var swap bool
	if dy > dx {
		dx, dy = dy, dx
		swap = true
	}

	err := 2*dy - dx

	for i := 0; i <= dx; i++ {
		if swap {
			ppm.Set(p1.Y, p1.X, color)
		} else {
			ppm.Set(p1.X, p1.Y, color)
		}

		for err >= 0 {
			if err > 0 || (err == 0 && xSign > 0) {
				p1.Y += ySign
			}
			err -= 2 * dx
		}

		if err < 0 || (err == 0 && xSign > 0) {
			p1.X += xSign
		}

		err += 2 * dy
	}
}

//DrawRectangle dessine un rectangle.
func (ppm *PPM) DrawRectangle(p1 Point, width, height int, color Pixel) {
	p2 := Point{p1.X + width, p1.Y}
	p3 := Point{p1.X + width, p1.Y + height}
	p4 := Point{p1.X, p1.Y + height}

	ppm.DrawLine(p1, p2, color)
	ppm.DrawLine(p2, p3, color)
	ppm.DrawLine(p3, p4, color)
	ppm.DrawLine(p4, p1, color)
}

// DrawFilledRectangle dessine un rectangle rempli.
func (ppm *PPM) DrawFilledRectangle(p1 Point, width, height int, color Pixel) {
	for y := p1.Y; y < p1.Y+height; y++ {
		for x := p1.X; x < p1.X+width; x++ {
			ppm.Set(x, y, color)
		}
	}
}

// DrawCircle dessine un cercle.
func (ppm *PPM) DrawCircle(center Point, radius int, color Pixel) {
	x := radius
	y := 0
	err := 0

	for x >= y {
		ppm.Set(center.X+x, center.Y+y, color)
		ppm.Set(center.X+y, center.Y+x, color)
		ppm.Set(center.X-y, center.Y+x, color)
		ppm.Set(center.X-x, center.Y+y, color)
		ppm.Set(center.X-x, center.Y-y, color)
		ppm.Set(center.X-y, center.Y-x, color)
		ppm.Set(center.X+y, center.Y-x, color)
		ppm.Set(center.X+x, center.Y-y, color)

		if err <= 0 {
			y += 1
			err += 2*y + 1
		}

		if err > 0 {
			x -= 1
			err -= 2*x + 1
		}
	}
}

// Draw Filled Circle dessine un cercle rempli.
func (ppm *PPM) DrawFilledCircle(center Point, radius int, color Pixel) {
	for y := -radius; y <= radius; y++ {
		for x := -radius; x <= radius; x++ {
			if x*x+y*y <= radius*radius {
				ppm.Set(center.X+x, center.Y+y, color)
			}
		}
	}
}

// DrawTriangle dessine un triangle.
func (ppm *PPM) DrawTriangle(p1, p2, p3 Point, color Pixel) {
	ppm.DrawLine(p1, p2, color)
	ppm.DrawLine(p2, p3, color)
	ppm.DrawLine(p3, p1, color)
}

// DrawFilledTriangle dessine un triangle rempli.
func (ppm *PPM) DrawFilledTriangle(p1, p2, p3 Point, color Pixel) {
	// Utilise le balayage de lignes pour remplir le triangle
	if p1.Y > p2.Y {
		p1, p2 = p2, p1
	}
	if p2.Y > p3.Y {
		p2, p3 = p3, p2
	}
	if p1.Y > p2.Y {
		p1, p2 = p2, p1
	}

	totalHeight := p3.Y - p1.Y
	for y := p1.Y; y <= p2.Y; y++ {
		segmentHeight := p2.Y - p1.Y + 1
		alpha := float64(y-p1.Y) / float64(totalHeight)
		beta := float64(y-p1.Y) / float64(segmentHeight)

		A := Point{int(float64(p1.X) + float64(p3.X-p1.X)*alpha), y}
		B := Point{int(float64(p1.X) + float64(p2.X-p1.X)*beta), y}

		ppm.DrawLine(A, B, color)
	}

	totalHeight = p3.Y - p2.Y
	for y := p2.Y; y <= p3.Y; y++ {
		segmentHeight := p3.Y - p2.Y + 1
		alpha := float64(y-p1.Y) / float64(totalHeight)
		beta := float64(y-p2.Y) / float64(segmentHeight)

		A := Point{int(float64(p1.X) + float64(p3.X-p1.X)*alpha), y}
		B := Point{int(float64(p2.X) + float64(p3.X-p2.X)*beta), y}

		ppm.DrawLine(A, B, color)
	}
}

// DrawPolygon dessine un polygone.
func (ppm *PPM) DrawPolygon(points []Point, color Pixel) {
	numPoints := len(points)
	for i := 0; i < numPoints-1; i++ {
		ppm.DrawLine(points[i], points[i+1], color)
	}
	ppm.DrawLine(points[numPoints-1], points[0], color)
}

// DrawPolygon dessine un polygone.
func (ppm *PPM) DrawPolygon(points []Point, color Pixel) {
	numPoints := len(points)
	for i := 0; i < numPoints-1; i++ {
		ppm.DrawLine(points[i], points[i+1], color)
	}
	ppm.DrawLine(points[numPoints-1], points[0], color)
}

// DrawFilledPolygon dessine un polygone rempli.
func (ppm *PPM) DrawFilledPolygon(points []Point, color Pixel) {
	// Utilise la balayage de lignes pour remplir le polygone
	minY, maxY := points[0].Y, points[0].Y
	for _, p := range points {
		if p.Y < minY {
			minY = p.Y
		}
		if p.Y > maxY {
			maxY = p.Y
		}
	}

	// Remplit chaque ligne entre minY et maxY
	for y := minY; y <= maxY; y++ {
		intersections := []int{}

		// Trouve les intersections avec chaque segment du polygone
		for i := 0; i < len(points); i++ {
			p1 := points[i]
			p2 := points[(i+1)%len(points)]

			if (p1.Y <= y && p2.Y > y) || (p2.Y <= y && p1.Y > y) {
				xIntersection := int(float64(p1.X) + (float64(y-p1.Y)/float64(p2.Y-p1.Y))*(float64(p2.X-p1.X)))
				intersections = append(intersections, xIntersection)
			}
		}

		// Trie les intersections par ordre croissant
		sort.Ints(intersections)

		// Remplit les pixels entre les intersections
		for i := 0; i < len(intersections); i += 2 {
			x1 := intersections[i]
			x2 := intersections[i+1]

			for x := x1; x <= x2; x++ {
				ppm.Set(x, y, color)
			}
		}
	}
}

// DrawKochSnowflake dessine un flocon de Koch.
func (ppm *PPM) DrawKochSnowflake(n int, start Point, width int, color Pixel) {
	if n <= 0 {
		return
	}

	// Calcule les points du triangle équilatéral
	height := int(float64(width) * math.Sqrt(3.0) / 2.0)
	p1 := start
	p2 := Point{start.X + width, start.Y}
	p3 := Point{start.X + width / 2, start.Y - height}

	// Calcule les points des segments du flocon de Koch
	p4 := Point{(2*p1.X + p3.X) / 3, (2*p1.Y + p3.Y) / 3}
	p5 := Point{(p1.X + 2*p3.X) / 3, (p1.Y + 2*p3.Y) / 3}
	p6 := Point{p4.X + (p5.X-p4.X)*0.5 - (p5.Y-p4.Y)*math.Sqrt(3.0)/2.0, p4.Y + (p5.X-p4.X)*math.Sqrt(3.0)/2.0 + (p5.Y-p4.Y)*0.5}

	// Dessine les segments du flocon de Koch
	ppm.DrawLine(p1, p2, color)
	ppm.DrawLine(p2, p3, color)
	ppm.DrawLine(p3, p1, color)
	ppm.DrawLine(p4, p5, color)

	// Appelle récursivement pour les segments supplémentaires
	ppm.DrawKochSnowflake(n-1, p1, width/3, color)
	ppm.DrawKochSnowflake(n-1, p4, width/3, color)
	ppm.DrawKochSnowflake(n-1, p5, width/3, color)
	ppm.DrawKochSnowflake(n-1, p2, width/3, color)
	ppm.DrawKochSnowflake(n-1, p6, width/3, color)
	ppm.DrawKochSnowflake(n-1, p3, width/3, color)
}

// DrawSierpinskiTriangle dessine un triangle de Sierpinski.
func (ppm *PPM) DrawSierpinskiTriangle(n int, start Point, width int, color Pixel) {
	if n <= 0 {
		return
	}

	// Calcule les points du triangle équilatéral
	height := int(float64(width) * math.Sqrt(3.0) / 2.0)
	p1 := start
	p2 := Point{start.X + width, start.Y}
	p3 := Point{start.X + width / 2, start.Y - height}

	// Calcule le milieu des côtés
	mid1 := Point{(p1.X + p2.X) / 2, (p1.Y + p2.Y) / 2}
	mid2 := Point{(p2.X + p3.X) / 2, (p2.Y + p3.Y) / 2}
	mid3 := Point{(p3.X + p1.X) / 2, (p3.Y + p1.Y) / 2}

	// Calcule le centre du triangle
	center := Point{(p1.X + p2.X + p3.X) / 3, (p1.Y + p2.Y + p3.Y) / 3}

	// Dessine le triangle intérieur
	ppm.DrawTriangle(mid1, mid2, mid3, color)

	// Appelle récursivement pour les trois triangles restants
	ppm.DrawSierpinskiTriangle(n-1, start, width/2, color)
	ppm.DrawSierpinskiTriangle(n-1, mid1, width/2, color)
	ppm.DrawSierpinskiTriangle(n-1, mid2, width/2, color)
	ppm.DrawSierpinskiTriangle(n-1, mid3, width/2, color)
	ppm.DrawSierpinskiTriangle(n-1, center, width/2, color)
}

// DrawPerlinNoise dessine du bruit de Perlin.
func (ppm *PPM) DrawPerlinNoise(color1 Pixel, color2 Pixel) {
	// Génère une grille de vecteurs aléatoires
	grid := make([][]Point, ppm.Height)
	for y := 0; y < ppm.Height; y++ {
		grid[y] = make([]Point, ppm.Width)
		for x := 0; x < ppm.Width; x++ {
			angle := 2.0 * math.Pi * rand.Float64()
			grid[y][x] = Point{int(math.Cos(angle)), int(math.Sin(angle))}
		}
	}

	// Dessine le bruit de Perlin
	for y := 0; y < ppm.Height; y++ {
		for x := 0; x < ppm.Width; x++ {
			// Interpolation bilinéaire pour obtenir la valeur de bruit
			u := float64(x) / float64(ppm.Width-1)
			v := float64(y) / float64(ppm.Height-1)
			ix := u * float64(ppm.Width-1)
			iy := v * float64(ppm.Height-1)
			iu := int(ix)
			iv := int(iy)

			du := ix - float64(iu)
			dv := iy - float64(iv)

			// Calcule les valeurs interpolées
			c00 := float64(grid[iv][iu].X)*(1-du) + float64(grid[iv][iu].Y)*(1-dv)
			c10 := float64(grid[iv][iu+1].X)*du + float64(grid[iv][iu+1].Y)*(1-dv)
			c01 := float64(grid[iv+1][iu].X)*(1-du) + float64(grid[iv+1][iu].Y)*dv
			c11 := float64(grid[iv+1][iu+1].X)*du + float64(grid[iv+1][iu+1].Y)*dv

			noiseValue := (1-u)*(1-v)*c00 + u*(1-v)*c10 + (1-u)*v*c01 + u*v*c11

			// Interpolation linéaire entre les deux couleurs en fonction de la valeur de bruit
			interpolatedColor := InterpolateColor(color1, color2, noiseValue)

			// Applique la couleur au pixel
			ppm.Set(x, y, interpolatedColor)
		}
	}
}

// InterpolateColor réalise une interpolation linéaire entre deux couleurs.
func InterpolateColor(color1 Pixel, color2 Pixel, t float64) Pixel {
	r := uint8(float64(color1.R)*(1-t) + float64(color2.R)*t)
	g := uint8(float64(color1.G)*(1-t) + float64(color2.G)*t)
	b := uint8(float64(color1.B)*(1-t) + float64(color2.B)*t)
	return Pixel{r, g, b}
}