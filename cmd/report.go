/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/beppune/adcgo/download"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

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

		fmt.Println(reportname)

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

		request.Header.Add("Authorization", "Basic cmV0ZVxtYW56b2dpOToxS3J1bTFyMQ==")
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
