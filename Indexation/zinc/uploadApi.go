package zinc

import (
	"Indexation/parser"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime/pprof"
	"sync"
	"time"
)

const uploadPath = "http://localhost:4080/api/emails/_doc"
const uploadPathBulk = "http://localhost:4080/api/_bulkv2"
const userU = "admin"
const passwordU = "Complexpass#123"

// Con pushQuic usar 20 workers
const Workers = 20

func PushMailsQuickAP() {
	fmt.Println("Usando QuickAP")
	f, err := os.Create("cpuQuickAP.prof")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	err = pprof.StartCPUProfile(f)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer pprof.StopCPUProfile()

	start := time.Now()
	fmt.Println("Inicio:", start)

	files, err := getAllFiles("/home/gyron/Desktop/EmailIndex/Indexation/maildir/")
	if err != nil {
		log.Fatalf("Error al obtener la lista de archivos: %v", err)
	}

	//fmt.Println(len(files))

	blockOfMails := (len(files) + Workers - 1) / Workers
	//fmt.Println(blockOfMails)
	Blocks := make([][]string, 0)

	for i := 0; i < len(files); i += blockOfMails {
		end := i + blockOfMails
		if end > len(files) {
			end = len(files)
		}
		Blocks = append(Blocks, files[i:end])
	}

	fileChan := make(chan []string, len(Blocks))

	for _, block := range Blocks {
		fileChan <- block
	}
	close(fileChan)

	var wg sync.WaitGroup

	for i := 0; i < Workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for block := range fileChan {
				UploadEmailsBulk(block)

			}
		}()
	}

	wg.Wait()

	end := time.Now()
	fmt.Println("Fin:", end)
	fmt.Println("Tiempo transcurrido:", end.Sub(start))

}
func PushMailsQuick() {
	fmt.Println("Usando Quick")
	f, err := os.Create("cpuQuick.prof")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	err = pprof.StartCPUProfile(f)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer pprof.StopCPUProfile()

	start := time.Now()
	fmt.Println("Inicio:", start)

	files, err := getAllFiles("/home/gyron/Desktop/EmailIndex/Indexation/maildir/")
	if err != nil {
		log.Fatalf("Error al obtener la lista de archivos: %v", err)
	}

	fmt.Println(len(files))

	fileChan := make(chan string, len(files))

	for _, file := range files {
		fileChan <- file
	}
	close(fileChan)

	var wg sync.WaitGroup

	for i := 0; i < Workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for file := range fileChan {
				err := processFile(file)
				if err != nil {
					log.Printf("Error al procesar el archivo %s: %v", file, err)
				}
			}
		}()
	}

	wg.Wait()

	end := time.Now()
	fmt.Println("Fin:", end)
	fmt.Println("Tiempo transcurrido:", end.Sub(start))

}

func PushMails() {
	f, err := os.Create("cpu.prof")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	err = pprof.StartCPUProfile(f)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer pprof.StopCPUProfile()

	start := time.Now()
	fmt.Println("Inicio:", start)

	files, err := getAllFiles("/home/gyron/Desktop/EmailIndex/Indexation/maildir/")
	if err != nil {
		log.Fatalf("Error al obtener la lista de archivos: %v", err)
	}

	fmt.Println(len(files))

	for i := 0; i < len(files); i++ {
		fmt.Println(i)
		err := processFile(files[i])
		if err != nil {
			log.Printf("Error al procesar el archivo %s: %v", files[0], err)
		}
	}

	end := time.Now()
	fmt.Println("Fin:", end)
	fmt.Println("Tiempo transcurrido:", end.Sub(start))

	// Llama a la función EmailFromFile con la ruta del archivo como argumento
	// emailObj, err := parser.EmailFromFile("/home/gyron/Desktop/EmailIndex/Indexation/maildir/allen-p/inbox/5.")
	// if err != nil {
	// 	log.Fatalf("Error al leer el archivo de correo: %v", err)
	// }

	// UploadEmails(emailObj)

	// Imprime el objeto de correo obtenido
	// fmt.Printf("Mensaje ID: %s\n", emailObj.MessageId)
	// fmt.Printf("Fecha: %s\n", emailObj.Date)
	// fmt.Printf("De: %s\n", emailObj.From)
	// fmt.Printf("Para: %s\n", emailObj.To)
	// fmt.Printf("Asunto: %s\n", emailObj.Subject)
	// fmt.Printf("Cuerpo: %s\n", emailObj.Body)

}

func UploadEmails(email *parser.Email) error {
	jsonBytes, err := json.Marshal(&email)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", uploadPath, bytes.NewReader(jsonBytes))
	if err != nil {
		return err
	}
	req.SetBasicAuth(userU, passwordU)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return fmt.Errorf("zinc server responded with code %v: %v", resp.StatusCode, string(body))
	}

	return nil
}

type BulkEmails struct {
	Index   string          `json:"index"`
	Records []*parser.Email `json:"records"`
}

func UploadEmailsBulk(emailsBlock []string) error {
	parsedEmailsAndConvertedEmails := make([]*parser.Email, 0)
	for _, emailPath := range emailsBlock {
		emailObj, _ := parser.EmailFromFile(emailPath)
		// jsonBytes, _ := json.Marshal(emailObj)
		// fmt.Println(string(jsonBytes))
		parsedEmailsAndConvertedEmails = append(parsedEmailsAndConvertedEmails, emailObj)
	}

	//fmt.Println(len(parsedEmailsAndConvertedEmails))

	bulk := BulkEmails{
		Index:   "emails",
		Records: parsedEmailsAndConvertedEmails,
	}

	queryJSON, err := json.Marshal(bulk)
	if err != nil {
		log.Printf("Error al convertir la consulta a JSON: %v", err)
	}
	//fmt.Println(string(queryJSON))

	req, err := http.NewRequest("POST", uploadPathBulk, bytes.NewReader(queryJSON))
	if err != nil {
		return fmt.Errorf("error armando la request: %e", err)
	}
	req.SetBasicAuth(userU, passwordU)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("error mandando la request: %e", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return fmt.Errorf("zinc server responded with code %v: %v", resp.StatusCode, string(body))
	}

	return nil

}

func getAllFiles(dirPath string) ([]string, error) {
	var files []string
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}

func processFile(filePath string) error {
	// Parsear los documentos
	emailObj, err := parser.EmailFromFile(filePath)
	if err != nil {
		return fmt.Errorf("error al leer el archivo de correo %s: %v", filePath, err)
	}

	// Sube el correo electrónico
	err = UploadEmails(emailObj)
	if err != nil {
		return fmt.Errorf("error al subir el correo electrónico %s: %v", emailObj.MessageId, err)
	}

	//fmt.Printf("%s\n", emailObj.MessageId)

	return nil
}
