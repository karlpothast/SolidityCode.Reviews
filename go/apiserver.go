package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/acme/autocert"
)

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.Use(setCSPHeaders)
	router.HandleFunc("/", index)
	router.HandleFunc("/uploadfile", uploadFile)
	router.HandleFunc("/analyzefile/{id}", analyzeFile)
	router.HandleFunc("/getcategories", getCategories)
	router.HandleFunc("/analyzebycontractid", analyzeByContractId)

	certManager := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist("web3api.tools", "web3api.tools"),
		Cache:      autocert.DirCache("certs"),
	}

	server := &http.Server{
		Addr:    ":https",
		Handler: router,
		TLSConfig: &tls.Config{
			GetCertificate: certManager.GetCertificate,
		},
	}

	http.HandleFunc("/", httpRequestHandler)
	go http.ListenAndServe(":http", certManager.HTTPHandler(nil))
	log.Fatal(server.ListenAndServeTLS("", ""))
}

func setCSPHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self'; object-src 'self';style-src 'self' img-src 'self'; media-src 'self'; frame-ancestors 'self'; frame-src 'self'; connect-src 'self'")
		next.ServeHTTP(w, r)
	})
}

func index(w http.ResponseWriter, r *http.Request) {
	response := GenericResponse{Message: "API Server running "}
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		fmt.Println(err)
	}
}

func httpRequestHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "http req handler")
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
	// Maximum upload of 10 MB
	r.ParseMultipartForm(10 << 20)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Get handler for filename, size and headers
	file, handler, err := r.FormFile("selectedFile")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}

	defer file.Close()
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	// Create file
	dstFilePath := "./files/" + handler.Filename
	dst, err := os.Create(dstFilePath)
	if err != nil {
		fmt.Printf("Error on os.Create ")
		fmt.Printf(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Copy the uploaded file to the filesystem
	if _, err := io.Copy(dst, file); err != nil {
		fmt.Printf("Error on io.Copy")
		fmt.Printf(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println("Executing shell script")
	AnalysisResults_stdout := execShellCommand_RunAnalysis(dstFilePath, false)
	AnalysisResults_json := execShellCommand_RunAnalysis(dstFilePath, true)
	execShellCommand_CleanDirectory(dstFilePath)

	fmt.Println("Encoding JSON response")
	response := SmartContractReview{FileProcessed: handler.Filename, AnalysisResults_stdout: AnalysisResults_stdout, AnalysisResults_json: AnalysisResults_json}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		fmt.Println("Encoding JSON response error")
		fmt.Println(err)
	}
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		uploadFile(w, r)
	}
}

func analyzeFile(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
        w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
        w.Header().Set("Access-Control-Allow-Origin", "*")

        var fileId string
        vars := mux.Vars(r)
        fileId = vars["id"]

	var filePath string

	switch fileId {
        case "1":
	  filePath = "example-files/parity-wallet-hack.sol"
	}

        fmt.Println("Executing shell script")
        AnalysisResults_stdout := execShellCommand_RunAnalysis(filePath, false)
        AnalysisResults_json := execShellCommand_RunAnalysis(filePath, true)

        fmt.Println("Encoding JSON response")
        response := SmartContractReview{FileProcessed: filePath, AnalysisResults_stdout: AnalysisResults_stdout, AnalysisResults_json: AnalysisResults_json}
        err := json.NewEncoder(w).Encode(response)
        if err != nil {
                fmt.Println("Encoding JSON response error")
                fmt.Println(err)
        }
}

func getCategories(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	fmt.Println("Executing shell script")
	CategoryList := execShellCommand_GetCategories()

	fmt.Println("Encoding JSON response")
	response := ShellOutput{Output: CategoryList}
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		fmt.Println("Encoding JSON response error")
		fmt.Println(err)
	}
}

func analyzeByContractId(w http.ResponseWriter, r *http.Request) {
	fmt.Println("form data")
	fmt.Println(r.FormValue("textAreaContractId"))
	param1 := r.FormValue("textAreaContractId")
        //param1 := r.URL.Query().Get("contractId")
        //fmt.Println("query string val")
        //fmt.Println(param1)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	contractId := param1;
	fmt.Println("Executing shell script")
	AnalysisResults_stdout := execShellCommand_AnalyzeByContractId(contractId, false)
	AnalysisResults_json := execShellCommand_AnalyzeByContractId(contractId, true)

	fmt.Println("Encoding JSON response")
	response := SmartContractReview{FileProcessed: "", AnalysisResults_stdout: AnalysisResults_stdout, AnalysisResults_json: AnalysisResults_json}
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		fmt.Println("Encoding JSON response error")
		fmt.Println(err)
	}
}

func execShellCommand_RunAnalysis(fullPath string, jsonFormat bool) string {
	scriptName := "analyze-smart-contract.sh"
	if jsonFormat {
		scriptName = "analyze-smart-contract-json.sh"
	}

	cmd := exec.Command("bash", scriptName, "-n", fullPath)
	out, err := cmd.CombinedOutput()
	if err != nil {
		//log.Fatal(err)
		fmt.Println("Shell script execute error")
		fmt.Println(err)
	}
	return string(out)
}

func execShellCommand_AnalyzeByContractId(contractId string, jsonFormat bool) string {
	scriptName := "analyze-by-contract-id.sh"
	if jsonFormat {
		scriptName = "analyze-by-contract-id-json.sh"
	}

	cmd := exec.Command("bash", scriptName, "-n", contractId)
	out, err := cmd.CombinedOutput()
	if err != nil {
		//log.Fatal(err)
		fmt.Println("analyze by contract id execute error")
		fmt.Println(err)
	}
	return string(out)
}

func execShellCommand_GetCategories() string {
	scriptName := "get-categories.sh"
	cmd := exec.Command("bash", scriptName)
	out, err := cmd.CombinedOutput()
	if err != nil {
		//log.Fatal(err)
		fmt.Println("get-categories script execute error")
		fmt.Println(err)
	}
	return string(out)
}

func execShellCommand_CleanDirectory(fullPath string) string {
	scriptName := "clean-directory.sh"
	cmd := exec.Command("bash", scriptName, "-n", fullPath)
	out, err := cmd.CombinedOutput()
	if err != nil {
		//log.Fatal(err)
		fmt.Println("cleanup script execute error")
		fmt.Println(err)
	}
	return string(out)
}

type SmartContractReview struct {
	FileProcessed          string `json:"file_processed"`
	AnalysisResults_stdout string `json:"analysis_results_stdout"`
	AnalysisResults_json   string `json:"analysis_results_json"`
}

type GenericResponse struct {
	Message string `json:"message"`
}

type ShellOutput struct {
	Output string `json:"output"`
}
