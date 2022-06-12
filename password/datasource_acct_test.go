package password

import (
	_ "embed"
	"fmt"
	"github.com/hashicorp/packer-plugin-sdk/acctest"
	"os/exec"
	"testing"
)

//go:embed test-fixtures/outputs.pkr.hcl
var testDatasourceTemplateOutputs string

func check(buildCommand *exec.Cmd, logfile string) error {
	if buildCommand.ProcessState != nil {
		exitStatus := buildCommand.ProcessState.ExitCode()
		if exitStatus != 0 {
			return fmt.Errorf("non-zero exit code: (%d) - check log in %s: ", exitStatus, logfile)
		}
	}
	return nil
}

func TestAccPassword(t *testing.T) {
	testGeneratePasswordOutputs := &acctest.PluginTestCase{
		Name:     "password_datasource_generate_outputs",
		Template: testDatasourceTemplateOutputs,
		Check:    check,
	}

	acctest.TestPlugin(t, testGeneratePasswordOutputs)
}
