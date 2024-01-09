package cmd

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateRulesetCommand(t *testing.T) {
	cmd := GetGenerateRulesetCommand()
	b := bytes.NewBufferString("")
	cmd.SetOut(b)
	cmd.SetArgs([]string{
		"all",
		"test-output",
	})
	cmdErr := cmd.Execute()
	outBytes, err := io.ReadAll(b)

	assert.NoError(t, cmdErr)
	assert.NoError(t, err)
	assert.NotNil(t, outBytes)
	defer os.Remove("test-output-all.yaml")
}

func TestGenerateRulesetCommand_Recommended(t *testing.T) {
	cmd := GetGenerateRulesetCommand()
	b := bytes.NewBufferString("")
	cmd.SetOut(b)
	cmd.SetArgs([]string{
		"recommended",
		"test-output",
	})
	cmdErr := cmd.Execute()
	outBytes, err := io.ReadAll(b)

	assert.NoError(t, cmdErr)
	assert.NoError(t, err)
	assert.NotNil(t, outBytes)
	defer os.Remove("test-output-recommended.yaml")
}

func TestGenerateRulesetCommand_InvalidType(t *testing.T) {
	cmd := GetGenerateRulesetCommand()
	b := bytes.NewBufferString("")
	cmd.SetOut(b)
	cmd.SetArgs([]string{
		"fish-cakes",
		"test-output",
	})
	cmdErr := cmd.Execute()
	assert.Error(t, cmdErr)
}

func TestGenerateRulesetCommand_NoArgs(t *testing.T) {
	cmd := GetGenerateRulesetCommand()
	b := bytes.NewBufferString("")
	cmd.SetOut(b)
	cmd.SetArgs([]string{})
	cmdErr := cmd.Execute()
	assert.Error(t, cmdErr)
}

func TestGenerateRulesetCommand_BadWrite(t *testing.T) {
	cmd := GetGenerateRulesetCommand()
	b := bytes.NewBufferString("")
	cmd.SetOut(b)
	cmd.SetArgs([]string{
		"recommended",
		"/no/no/no-stop-/",
	})
	cmdErr := cmd.Execute()
	assert.Error(t, cmdErr)
}
