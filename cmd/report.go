/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"embed"

	"github.com/beppune/adcgo/download"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/term"
)

//go:embed body.txt
var bodycontent embed.FS

//go:embed report_template.xlsx
var xlstemplate embed.FS

func credentials() (string, string, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter Username: ")
	username, err := reader.ReadString('\n')
	if err != nil {
		return "", "", err
	}

	fmt.Print("Enter Password: ")
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", "", err
	}

	password := string(bytePassword)
	return strings.TrimSpace(username), strings.TrimSpace(password), nil
}

// reportCmd represents the report command
var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Download and produce reports from different sources",
	Long:  `A longer description`,

	Run: func(cmd *cobra.Command, args []string) {

		rawurl := viper.GetString("report.dailyurl")
		date, _ := cmd.PersistentFlags().GetString("date")
		t, _ := time.Parse(`2006-01-02`, date)
		date = t.Format(`02 01 2006`)
		reportname := viper.GetString("report.format")
		reportname = fmt.Sprintf(reportname, t.Format(`02_01_2006`), time.Now().Format(`02012006030405`))

		noclean := viper.GetBool("report.noclean")

		username, password, err := credentials()
		if err != nil {
			panic(err.Error())
		}

		username = `rete\` + username
		token := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))

		if rawurl == "" {
			panic("dailyurl required")
		}

		data, _ := bodycontent.ReadFile("body.txt")
		body := download.PrepareBody(date, data)

		request, err := http.NewRequest("POST", rawurl, body)
		if err != nil {
			panic(err.Error())
		}

		download.Prepare(request, rawurl)

		request.Header.Add("Authorization", "Basic "+token)
		request.Header.Add("Coockie", `ASP.NET_SessionId=jkdouw23z3q1itn0hu0mhs03`)

		client := &http.Client{}

		res, err := client.Do(request)
		if err != nil {
			panic(err.Error())
		}

		records, _ := download.ParseTable(res.Body)

		template, _ := xlstemplate.Open("report_template.xlsx")

		download.ProduceExcel(template, records, reportname)

		if !noclean {
			os.Remove("temp.txt")
		}

		if d, _ := cmd.PersistentFlags().GetBool("open"); d {
			excelPath := viper.GetString("excel-path")
			c := exec.Command(excelPath, reportname)
			err = c.Start()
			if err != nil {
				fmt.Println(err.Error())
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(reportCmd)

	reportCmd.PersistentFlags().String("dailyurl", "http://10.194.137.36/ACCESSIDC/ReportGiornaliero.aspx", "Daily report url")
	viper.BindPFlag("report.dailyurl", reportCmd.PersistentFlags().Lookup("dailyurl"))

	viper.Set("report.format", `ReportGiornaliero_TO1__%s_%s.xlsx`)

	reportCmd.PersistentFlags().Bool("noclean", false, "Do not clean temporary files")
	viper.BindPFlag("report.noclean", reportCmd.PersistentFlags().Lookup("noclean"))

	reportCmd.PersistentFlags().Bool("open", false, "Open generated report into excel")

	today := time.Now().Format("2006-01-02")
	reportCmd.PersistentFlags().String("date", today, "Report date (default today)")
}
