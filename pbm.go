package netpbm

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// PBM est une structure pour représenter des images PBM.
type PBM struct {
	Data          [][]bool
	Width, Height int
	MagicNumber   string
}

// ReadPBM lit une image PBM à partir d'un fichier et renvoie une structure représentant l'image.
func ReadPBM(filename string) (*PBM, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// Lire le numéro magique
	if !scanner.Scan() {
		return nil, fmt.Errorf("échec de la lecture du numéro magique")
	}
	magicNumber := scanner.Text()

	if magicNumber != "P1" && magicNumber != "P4" {
		return nil, fmt.Errorf("format PBM non pris en charge : %s", magicNumber)
	}

	// Ignorer les commentaires et les lignes vides
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) > 0 && line[0] != '#' {
			break
		}
	}

	// Lire la largeur et la hauteur
	if scanner.Err() != nil {
		return nil, fmt.Errorf("erreur lors de la lecture de la ligne de dimensions : %v", scanner.Err())
	}
	dimensions := strings.Fields(scanner.Text())
	if len(dimensions) != 2 {
		return nil, fmt.Errorf("ligne de dimensions invalide")
	}

	width, err := strconv.Atoi(dimensions[0])
	if err != nil {
		return nil, fmt.Errorf("échec de l'analyse de la largeur : %v", err)
	}

	height, err := strconv.Atoi(dimensions[1])
	if err != nil {
		return nil, fmt.Errorf("échec de l'analyse de la hauteur : %v", err)
	}

	// lire les données
	var data [][]bool
	for scanner.Scan() {
		line := scanner.Text()
		tokens := strings.Fields(line)
		row := make([]bool, width)
		for i, token := range tokens {
			if i >= width {
				break
			}
			if token == "1" {
				row[i] = true
			} else if token == "0" {
				row[i] = false
			} else {
				return nil, fmt.Errorf("caractère non valide dans les données : %s", token)
			}
		}
		data = append(data, row)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("erreur lors de la lecture du fichier : %v", err)
	}

	return &PBM{
		Data:        data,
		Width:       width,
		Height:      height,
		MagicNumber: magicNumber,
	}, nil
}

// Size retourne la largeur et la hauteur de l'image
func (pbm *PBM) Size() (int,int){
   return pbm.Width, pbm.Height // width = largeur ; height = hauteur (de l'image)

}

// At retourne la valeur du pixel a (x, y).
func (pbm *PBM) At(x, y int) bool{
	 // Vérifie si les coordonnées sont correctes
	if x < 0 || x >= pbm.Width || y < 0 || y >= pbm.Height {
		// Coordonnées invalides, retourne une valeur par défaut ou gére l'erreur
		return false
	}

	// Récupére la valeur du pixel aux coordonnées (x, y)
	return pbm.Data[y][x]
}

func (pbm *PBM) Set(x, y int, value bool) {
    // Vérifier si les coordonnées sont valides
    if x >= 0 && x < pbm.Width && y >= 0 && y < pbm.Height {
        // Modifier la valeur du pixel aux coordonnées (x, y)
        pbm.Data[y][x] = value
    }
    // Si les coordonnées sont invalides, ne rien faire
}

func (pbm *PBM) Save(filename string) error {
    // Créer un nom de fichier unique avec la date et l'heure du fichier
    horodatage := time.Now().Format("2006-01-02-15-04") // exemple de format
    newFichier := fmt.Sprintf("%s%s", strings.TrimSuffix(filename, ".pbm"), horodatage)

    // Créer ou ouvrir le fichier
    fichier, err := os.Create(newFichier)
    if err != nil {
        return fmt.Errorf("échec de la création du fichier : %v", err)
    }
    defer fichier.Close()

    // Écrire le numéro magique
    _, err = fmt.Fprintf(fichier, "%s\n", pbm.MagicNumber)
    if err != nil {
        return fmt.Errorf("échec de l'écriture du numéro magique : %v", err)
    }

    // Écrire les dimensions
    _, err = fmt.Fprintf(fichier, "%d %d\n", pbm.Width, pbm.Height)
    if err != nil {
        return fmt.Errorf("échec de l'écriture des dimensions : %v", err)
    }

    // Écrire les données
	for _, ligne := range pbm.Data {
		for _, pixel := range ligne {
			var valeur int
			if pixel {
				valeur = '∎'
			} else {
				valeur = '□'
			}
			_, err := fmt.Fprintf(fichier, "%d ", valeur)
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
func (pbm *PBM) Invert() {
    for y := 0; y < pbm.Height; y++ {
        for x := 0; x < pbm.Width; x++ {
            // inverse la valeur de chaque pixel
            pbm.Data[y][x] = !pbm.Data[y][x]
        }
    }

	// Sauvegarde l'image inversée
	bufio.ErrBufferFull = pbm.Save("image_inverse.pbm")
	if bufio.ErrBufferFull != nil {
   		fmt.Println("Erreur lors de la sauvegarde de l'image inversée :", bufio.ErrBufferFull)
   		return
	}
}

func (pbm *PBM) FlipAndFlop() {
    // Inverser horizontalement (flip)
    for y := 0; y < pbm.Height; y++ {
        for x := 0; x < pbm.Width/2; x++ {
            // Échanger les pixels symétriques horizontalement
            pbm.Data[y][x], pbm.Data[y][pbm.Width-x-1] = pbm.Data[y][pbm.Width-x-1], pbm.Data[y][x]
        }
    }

    // Inverser verticalement (flop)
    for y := 0; y < pbm.Height/2; y++ {
        for x := 0; x < pbm.Width; x++ {
            // Échanger les lignes symétriques verticalement
            pbm.Data[y][x], pbm.Data[pbm.Height-y-1][x] = pbm.Data[pbm.Height-y-1][x], pbm.Data[y][x]
        }
    }
}

func (pbm *PBM) SetMagicNumber(magicNumber string) {
    pbm.MagicNumber = magicNumber
}
