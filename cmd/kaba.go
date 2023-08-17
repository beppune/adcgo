/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io"
	"net"
	"os"
	"time"

	"github.com/hirochachacha/go-smb2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// kabaCmd represents the kaba command
var kabaCmd = &cobra.Command{
	Use:   "kaba",
	Short: "A brief description of your command",
	Long:  ` `,
	Run: func(cmd *cobra.Command, args []string) {
		server := viper.GetString("kaba.server")
		sharename := viper.GetString("kaba.sharename")
		root := viper.GetString("kaba.root")
		month, _ := cmd.PersistentFlags().GetString("month")
		date, _ := cmd.PersistentFlags().GetString("date")

		username, password, err := credentials()
		if err != nil {
			panic(err)
		}

		conn, err := net.Dial("tcp", server+":445")
		if err != nil {
			panic(err)
		}
		defer conn.Close()

		dialer := &smb2.Dialer{
			Initiator: &smb2.NTLMInitiator{
				User:     username,
				Password: password,
			},
		}

		session, err := dialer.Dial(conn)
		if err != nil {
			panic(err)
		}
		defer session.Logoff()

		fs, err := session.Mount(sharename)
		if err != nil {
			panic(err)
		}
		defer fs.Umount()

		path := root + "\\" + year + "\\" + month + "\\Allarmi\\" + date + month + "_ALL.xlsx"

		downloadFile(fs, path, date+month+"_ALL.xlsx")

		path = root + "\\" + year + "\\" + month + "\\Passaggi\\" + date + month + ".xlsx"

		downloadFile(fs, path, date+month+".xlsx")
	},
}

var year string

func downloadFile(fs *smb2.Share, path string, output string) {
	f, err := fs.Open(path)
	if err != nil {
		panic(err)
	}

	b, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(output, b, 0644)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func init() {
	rootCmd.AddCommand(kabaCmd)

	kabaCmd.PersistentFlags().String("server", "rto1y11c013", "smb server")
	viper.BindPFlag("kaba.server", kabaCmd.PersistentFlags().Lookup("server"))

	kabaCmd.PersistentFlags().String("sharename", "Dati", "sharename to mount")
	viper.BindPFlag("kaba.sharename", kabaCmd.PersistentFlags().Lookup("sharename"))

	kabaCmd.PersistentFlags().String("root", "DC\\REPORT\\Export\\Exos", "exports root directory of given share")
	viper.BindPFlag("kaba.root", kabaCmd.PersistentFlags().Lookup("root"))

	t := time.Now().Add(-(time.Hour * 24))
	year = t.Format("2006")

	kabaCmd.PersistentFlags().String("date", t.Format("02"), "date")
	kabaCmd.PersistentFlags().String("month", t.Format("01"), "month")
}
