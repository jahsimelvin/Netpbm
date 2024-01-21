package netpbm

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// PBM is a structure to represent PBM images.
type PGM struct {
	Data          [][]uint8
	Width, Height int
	MagicNumber   string
	Max uint
}

type Pixel struct{
	R, G, B uint8
}

// ReadPBM reads a PBM image from a file and returns a struct that represents the image.
func ReadPGM(filename string) (*PGM, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// lire le nombre magique
	if !scanner.Scan() {
		return nil, fmt.Errorf("failed to read magic number")
	}
	magicNumber := scanner.Text()

	if magicNumber != "P1" && magicNumber != "P4" {
		return nil, fmt.Errorf("unsupported PGM format: %s", magicNumber)
	}

	// Skip comments and empty lines
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) > 0 && line[0] != '#' {
			break
		}
	}

	// Read width and height
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

	// Read data
	var data [][]uint8
	for scanner.Scan() {
		line := scanner.Text()
		tokens := strings.Fields(line)
		row := make([]uint8, width)
		for i, token := range tokens {
			if i >= width {
				break
			}
			value, err := strconv.ParseUint(token, 10, 8)
			if err != nil {
				return nil, fmt.Errorf("caractère non valide dans les données : %s", token)
			}
			row[i] = uint8(value)
		}
		data = append(data, row)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	return &PGM{
		Data:        data,
		Width:       width,
		Height:      height,
		MagicNumber: magicNumber,
	}, nil
}

// Size retourne la largeur et la hauteur de l'image
func (pgm *PGM) Size() (int,int){
	return pgm.Width, pgm.Height // width = largeur ; height = hauteur (de l'image)

}

// At retourne la valeur du pixel a (x, y).
func (pgm *PGM) At(x, y int) uint8{
	// Vérifie si les coordonnées sont correctes
   if x < 0 || x >= pgm.Width || y < 0 || y >= pgm.Height {
	   // Coordonnées invalides, retourne une valeur par défaut ou gére l'erreur
	   return 0
   }

   // Récupére la valeur du pixel aux coordonnées (x, y)
   return pgm.Data[y][x]
}

func (pgm *PGM) Set(x, y int, value uint8) {
    // Vérifier si les coordonnées sont valides
    if x >= 0 && x < pgm.Width && y >= 0 && y < pgm.Height {
        // Modifier la valeur du pixel aux coordonnées (x, y)
        pgm.Data[y][x] = value
    }
    // Si les coordonnées sont invalides, ne rien faire
}

func (pgm *PGM) Save(filename string) error {
    // Créer un nom de fichier unique avec un horodatage
    horodatage := time.Now().Format("2006-01-02-15-04")
    nouveauNomFichier := fmt.Sprintf("%s%s", strings.TrimSuffix(filename, ".pbm"), horodatage)

    // Créer ou ouvrir le fichier
    fichier, err := os.Create(nouveauNomFichier)
    if err != nil {
        return fmt.Errorf("échec de la création du fichier : %v", err)
    }
    defer fichier.Close()

    // Écrire le numéro magique
    _, err = fmt.Fprintf(fichier, "%s\n", pgm.MagicNumber)
    if err != nil {
        return fmt.Errorf("échec de l'écriture du numéro magique : %v", err)
    }

    // Écrire les dimensions
    _, err = fmt.Fprintf(fichier, "%d %d\n", pgm.Width, pgm.Height)
    if err != nil {
        return fmt.Errorf("échec de l'écriture des dimensions : %v", err)
    }

    // Écrire les données
    for _, ligne := range pgm.Data {
        for _, pixel := range ligne {
            _, err := fmt.Fprintf(fichier, "%d ", pixel)
            if err != nil {
                return fmt.Errorf("échec de l'écriture de la valeur du pixel : %v", err)
            }
        }
        _, err := fmt.Fprintln(fichier) // Nouvelle ligne après chaque ligne de pixels
        if err != nil {
            return fmt.Errorf("échec de l'écriture d'un saut de ligne : %v", err)
        }
    }

    return nil
}

// inverse les couleurs de chaque pixel de l'image pbm
func (pgm *PGM) Invert() {
    for y := 0; y < pgm.Height; y++ {
        for x := 0; x < pgm.Width; x++ {
            // inverse la valeur de chaque pixel
            pgm.Data[y][x] = uint8(pgm.Max) - pgm.Data[y][x] + 1
        }
    }
}

func (pgm *PGM) FlipAndFlop() {
    // Inverser horizontalement (flip)
    for y := 0; y < pgm.Height; y++ {
        for x := 0; x < pgm.Width/2; x++ {
            // Échanger les pixels symétriques horizontalement
            pgm.Data[y][x], pgm.Data[y][pgm.Width-x-1] = pgm.Data[y][pgm.Width-x-1], pgm.Data[y][x]
        }
    }

    // Inverser verticalement (flop)
    for y := 0; y < pgm.Height/2; y++ {
        for x := 0; x < pgm.Width; x++ {
            // Échanger les lignes symétriques verticalement
            pgm.Data[y][x], pgm.Data[pgm.Height-y-1][x] = pgm.Data[pgm.Height-y-1][x], pgm.Data[y][x]
        }
    }
}

func (pgm *PGM) SetMagicNumber(magicNumber string) {
    pgm.MagicNumber = magicNumber
}

// SetMaxValue sets the max value of the PGM image.
func (pgm *PGM) SetMaxValue(maxValue uint8) {
	pgm.Max = uint(maxValue)
}

// Rotate90CW rotates the PGM image 90° clockwise.
func (pgm *PGM) Rotate90CW() {
	// Transpose the matrix
	for i := 0; i < pgm.Height; i++ {
		for j := i + 1; j < pgm.Width; j++ {
			pgm.Data[i][j], pgm.Data[j][i] = pgm.Data[j][i], pgm.Data[i][j]
		}
	}

	// Reverse each row
	for i := 0; i < pgm.Height; i++ {
		for j, k := 0, pgm.Width-1; j < k; j, k = j+1, k-1 {
			pgm.Data[i][j], pgm.Data[i][k] = pgm.Data[i][k], pgm.Data[i][j]
		}
	}
}

// ToPBM converts the PGM image to PBM.
func (pgm *PGM) ToPBM() *PBM {
    // Créer une nouvelle structure PBM
    pbm := &PBM{
        Width:       pgm.Width,
        Height:      pgm.Height,
        MagicNumber: "P1", // PBM a le numéro magique "P1"
    }

    // Initialiser les données de l'image PBM
    pbm.Data = make([][]uint8, pgm.Height)
    for y := 0; y < pgm.Height; y++ {
        pbm.Data[y] = make([]uint8, pgm.Width)
        for x := 0; x < pgm.Width; x++ {
            // Définir un seuil pour convertir en 1 ou 0
            if pgm.Data[y][x] > pgm.Max/2 {
                pbm.Data[y][x] = 1
            } else {
                pbm.Data[y][x] = 0
            }
        }
    }

    return pbm
}