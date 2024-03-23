package main

import (
	"BogdanFloris/langrunner/internal/spec"
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

var specs *spec.Spec

func runHandler(c *gin.Context) {
	var req RunRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	stdout, err := run(req.Code, req.Language)
	if err != nil {
		c.JSON(http.StatusOK, RunResponse{Type: "error", Payload: err.Error()})
		return
	}
	c.JSON(http.StatusOK, RunResponse{Type: "output", Payload: string(stdout)})
}

func run(code string, language string) ([]byte, error) {
	// Get the language spec
	langSpec, err := specs.Get(language)
	if err != nil {
		return nil, errors.New("invalid language")
	}

	// Create a temporary dir to store the code
	tmpDir, err := os.MkdirTemp("", fmt.Sprintf("%s-", language))
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tmpDir)
	codeFile, err := os.Create(fmt.Sprintf("%s/%s", tmpDir, langSpec.FileName))
	if err != nil {
		return nil, err
	}

	// Write the code to the file
	_, err = codeFile.WriteString(code)
	if err != nil {
		return nil, err
	}

	// Compile and run the code
	cmd := exec.Command(langSpec.GetCommandWithArgs()[0], langSpec.GetCommandWithArgs()[1:]...)
	cmd.Dir = tmpDir
	var stderr bytes.Buffer
	var stdout bytes.Buffer
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout
	err = cmd.Run()
	if err != nil {
		println(err.Error())
		return nil, errors.New(stderr.String())
	}
	return stdout.Bytes(), nil
}

func main() {
	var err error
	specs, err = spec.New("./spec/spec.toml")
	if err != nil {
		panic(err)
	}
	r := gin.Default()
	r.Use(cors.Default())
	r.POST("/run", runHandler)
	r.Run()
}
