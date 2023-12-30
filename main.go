package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type RunRequest struct {
	Language string `json:"language" binding:"required"`
	Code     string `json:"code"     binding:"required"`
}

type PayloadType string

const (
	PayloadTypeOutput PayloadType = "output"
	PayloadTypeError  PayloadType = "error"
)

type RunResponse struct {
	Type    PayloadType `json:"type"`
	Payload string      `json:"payload"`
}

func runHandler(c *gin.Context) {
	var req RunRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var stdout []byte
	var err error
	switch req.Language {
	case "rust":
		stdout, err = runRust(req.Code)
	case "java":
		stdout, err = runJava(req.Code)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid language"})
	}
	if err != nil {
		c.JSON(http.StatusOK, RunResponse{Type: "error", Payload: err.Error()})
		return
	}
	c.JSON(http.StatusOK, RunResponse{Type: "output", Payload: string(stdout)})
}

func runRust(code string) ([]byte, error) {
	// Create a temporary dir to store the Rust code
	tmpDir, err := os.MkdirTemp("", "rust-")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tmpDir)

	codeFile, err := os.Create(fmt.Sprintf("%s/main.rs", tmpDir))
	if err != nil {
		return nil, err
	}
	defer codeFile.Close()

	// Write the code to the file
	_, err = codeFile.WriteString(code)
	if err != nil {
		return nil, err
	}

	// Create output file
	outFile, err := os.Create(fmt.Sprintf("%s/main", tmpDir))
	if err != nil {
		return nil, err
	}
	defer outFile.Close()

	// Compile the code
	cmd := exec.Command(
		"rustc",
		codeFile.Name(),
		"-o",
		outFile.Name(),
	)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		return nil, errors.New(stderr.String())
	}

	// Run the compiled code and return the output
	cmd = exec.Command(outFile.Name())
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		println(err.Error())
		return nil, errors.New(stderr.String())
	}
	return stdout.Bytes(), nil
}

func runJava(code string) ([]byte, error) {
	// Create a temporary dir to store the Rust code
	tmpDir, err := os.MkdirTemp("", "java-")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tmpDir)

	codeFile, err := os.Create(fmt.Sprintf("%s/Main.java", tmpDir))
	if err != nil {
		return nil, err
	}
	defer codeFile.Close()

	// Write the code to the file
	_, err = codeFile.WriteString(code)
	if err != nil {
		return nil, err
	}

	// Create output file

	// Compile the code
	cmd := exec.Command(
		"javac",
		"-d",
		tmpDir,
		codeFile.Name(),
	)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		return nil, errors.New(stderr.String())
	}

	// Run the compiled code and return the output
	cmd = exec.Command("java", "Main")
	cmd.Dir = tmpDir
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		println(err.Error())
		return nil, errors.New(stderr.String())
	}
	return stdout.Bytes(), nil
}

func main() {
	r := gin.Default()
	r.Use(cors.Default())
	r.POST("/run", runHandler)
	r.Run()
}
