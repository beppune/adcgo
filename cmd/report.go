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
	"strings"
	"syscall"
	"time"

	"github.com/beppune/adcgo/download"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/term"
)

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
		template := viper.GetString("report.dailytemplate")
		bodyfile := viper.GetString("report.bodyfile")
		date := viper.GetString("date")
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

		if template == "" {
			panic("template required")
		}

		body := download.PrepareBody(date, bodyfile)

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

		//b, _ := io.ReadAll(res.Body)
		//os.WriteFile("res.dump.txt", b, 0644)

		records, _ := download.ParseTable(res.Body)

		download.ProduceExcel("report_template.xlsx", records, reportname)

		if !noclean {
			os.Remove("temp.txt")
		}
	},
}

func init() {
	rootCmd.AddCommand(reportCmd)

	reportCmd.PersistentFlags().String("dailyurl", "", "Daily report url")
	viper.BindPFlag("report.dailyurl", reportCmd.PersistentFlags().Lookup("dailyurl"))

	reportCmd.PersistentFlags().String("dailytemplate", "", "Daily report xlsx file")
	viper.BindPFlag("report.dailytemplate", reportCmd.PersistentFlags().Lookup("dailytemplate"))

	reportCmd.PersistentFlags().String("bodyfile", "body.txt", "Request body template")
	viper.BindPFlag("report.bodyfile", reportCmd.PersistentFlags().Lookup("bodyfile"))

	reportCmd.PersistentFlags().String("format", `ReportGiornaliero_TO1__%s_%s.xls`, "filename format for new export. Accept time.Format semantics")
	viper.BindPFlag("report.format", reportCmd.PersistentFlags().Lookup("format"))

	reportCmd.PersistentFlags().Bool("noclean", false, "Do not clean temporary files")
	viper.BindPFlag("report.noclean", reportCmd.PersistentFlags().Lookup("noclean"))

	today := time.Now().Format("2006-01-02")
	reportCmd.PersistentFlags().String("date", today, "Report date (default today)")
	viper.BindPFlag("date", reportCmd.PersistentFlags().Lookup("date"))
}
