package testsuites

import (
	"fmt"
	"os"
	"strconv"

	"github.com/robfig/cron"
	"github.com/spf13/cobra"

	"github.com/kubeshop/testkube/cmd/kubectl-testkube/commands/common"
	apiClient "github.com/kubeshop/testkube/pkg/api/v1/client"
	"github.com/kubeshop/testkube/pkg/api/v1/testkube"
	"github.com/kubeshop/testkube/pkg/crd"
	"github.com/kubeshop/testkube/pkg/ui"
)

func NewCreateTestSuitesCmd() *cobra.Command {

	var (
		name                     string
		file                     string
		labels                   map[string]string
		variables                []string
		secretVariables          []string
		schedule                 string
		executionName            string
		httpProxy, httpsProxy    string
		secretVariableReferences map[string]string
		timeout                  int32
		jobTemplate              string
		cronJobTemplate          string
		scraperTemplate          string
		pvcTemplate              string
		jobTemplateReference     string
		cronJobTemplateReference string
		scraperTemplateReference string
		pvcTemplateReference     string
		update                   bool
	)

	cmd := &cobra.Command{
		Use:     "testsuite",
		Aliases: []string{"testsuites", "ts"},
		Short:   "Create new TestSuite",
		Long:    `Create new TestSuite Custom Resource`,
		Run: func(cmd *cobra.Command, args []string) {
			crdOnly, err := strconv.ParseBool(cmd.Flag("crd-only").Value.String())
			ui.ExitOnError("parsing flag value", err)

			options, err := NewTestSuiteUpsertOptionsFromFlags(cmd)
			ui.ExitOnError("getting test suite options", err)

			if options.Name == "" {
				ui.Failf("pass valid test suite name (in '--name' flag)")
			}

			if !crdOnly {
				client, namespace, err := common.GetClient(cmd)
				ui.ExitOnError("getting client", err)

				testSuite, _ := client.GetTestSuite(options.Name)

				if options.Name == testSuite.Name {
					if cmd.Flag("update").Changed {
						if !update {
							ui.Failf("TestSuite with name '%s' already exists in namespace %s, ", testSuite.Name, namespace)
						}
					} else {
						var ok bool
						if stat, _ := os.Stdin.Stat(); (stat.Mode() & os.ModeCharDevice) != 0 {
							ok = ui.Confirm(fmt.Sprintf("TestSuite with name '%s' already exists in namespace %s, ", testSuite.Name, namespace) +
								"do you want to overwrite it?")
						}

						if !ok {
							ui.Failf("TestSuite creation was aborted")
						}
					}

					options, err := NewTestSuiteUpdateOptionsFromFlags(cmd)
					ui.ExitOnError("getting test suite options", err)

					testSuite, err = client.UpdateTestSuite(options)
					ui.ExitOnError("updating TestSuite "+testSuite.Name+" in namespace "+namespace, err)

					ui.SuccessAndExit("TestSuite updated", testSuite.Name)
				}

				_, err = client.CreateTestSuite(apiClient.UpsertTestSuiteOptions(options))
				ui.ExitOnError("creating test suite "+options.Name+" in namespace "+options.Namespace, err)

				ui.Success("Test suite created", options.Name)
			} else {
				(*testkube.TestSuiteUpsertRequest)(&options).QuoteTestSuiteTextFields()
				data, err := crd.ExecuteTemplate(crd.TemplateTestSuite, options)
				ui.ExitOnError("executing crd template", err)

				ui.Info(data)
			}
		},
	}

	cmd.Flags().StringVarP(&file, "file", "f", "", "JSON test suite file - will be read from stdin if not specified, look at testkube.TestUpsertRequest")
	cmd.Flags().StringVar(&name, "name", "", "Set/Override test suite name")
	cmd.Flags().StringToStringVarP(&labels, "label", "l", nil, "label key value pair: --label key1=value1")
	cmd.Flags().StringArrayVarP(&variables, "variable", "v", nil, "param key value pair: --variable key1=value1")
	cmd.Flags().StringArrayVarP(&secretVariables, "secret-variable", "s", nil, "secret variable key value pair: --secret-variable key1=value1")
	cmd.Flags().StringVarP(&schedule, "schedule", "", "", "test suite schedule in a cron job form: * * * * *")
	cmd.Flags().StringVarP(&executionName, "execution-name", "", "", "execution name, if empty will be autogenerated")
	cmd.Flags().StringVar(&httpProxy, "http-proxy", "", "http proxy for executor containers")
	cmd.Flags().StringVar(&httpsProxy, "https-proxy", "", "https proxy for executor containers")
	cmd.Flags().StringToStringVarP(&secretVariableReferences, "secret-variable-reference", "", nil, "secret variable references in a form name1=secret_name1=secret_key1")
	cmd.Flags().Int32Var(&timeout, "timeout", 0, "duration in seconds for test suite to timeout. 0 disables timeout.")
	cmd.Flags().StringVar(&jobTemplate, "job-template", "", "job template file path for extensions to job template")
	cmd.Flags().StringVar(&cronJobTemplate, "cronjob-template", "", "cron job template file path for extensions to cron job template")
	cmd.Flags().StringVar(&scraperTemplate, "scraper-template", "", "scraper template file path for extensions to scraper template")
	cmd.Flags().StringVar(&pvcTemplate, "pvc-template", "", "pvc template file path for extensions to pvc template")
	cmd.Flags().StringVar(&jobTemplateReference, "job-template-reference", "", "reference to job template to use for the test")
	cmd.Flags().StringVar(&cronJobTemplateReference, "cronjob-template-reference", "", "reference to cron job template to use for the test")
	cmd.Flags().StringVar(&scraperTemplateReference, "scraper-template-reference", "", "reference to scraper template to use for the test")
	cmd.Flags().StringVar(&pvcTemplateReference, "pvc-template-reference", "", "reference to pvc template to use for the test")
	cmd.Flags().BoolVar(&update, "update", false, "update, if test suite already exists")
	cmd.Flags().MarkDeprecated("disable-webhooks", "disable-webhooks is deprecated")
	cmd.Flags().MarkDeprecated("enable-webhooks", "enable-webhooks is deprecated")

	return cmd
}

func validateSchedule(schedule string) error {
	if schedule == "" {
		return nil
	}

	specParser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	if _, err := specParser.Parse(schedule); err != nil {
		return err
	}

	return nil
}
