package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/NebulousLabs/Sia/api"
	"github.com/NebulousLabs/Sia/modules"
)

var (
	renterCmd = &cobra.Command{
		Use:   "renter",
		Short: "Perform renter actions",
		Long:  "Upload and download files, or view a list of previously uploaded files.",
		Run:   wrap(renterstatuscmd),
	}

	renterUploadCmd = &cobra.Command{
		Use:   "upload [filename] [nickname]",
		Short: "Upload a file",
		Long:  "Upload a file using a given nickname.",
		Run:   wrap(renteruploadcmd),
	}

	renterDownloadCmd = &cobra.Command{
		Use:   "download [nickname] [destination]",
		Short: "Download a file",
		Long:  "Download a previously-uploaded file to a specified destination.",
		Run:   wrap(renterdownloadcmd),
	}

	renterDownloadQueueCmd = &cobra.Command{
		Use:   "queue",
		Short: "View the download queue",
		Long:  "View the list of files that have been downloaded.",
		Run:   wrap(renterdownloadqueuecmd),
	}

	renterStatusCmd = &cobra.Command{
		Use:   "status",
		Short: "View a list of uploaded files",
		Long:  "View a list of files that have been uploaded to the network.",
		Run:   wrap(renterstatuscmd),
	}
)

func renteruploadcmd(source, nickname string) {
	err := callAPI(fmt.Sprintf("/renter/upload?source=%s&nickname=%s", source, nickname))
	if err != nil {
		fmt.Println("Could not upload file:", err)
		return
	}
	fmt.Println("Upload initiated.")
}

func renterdownloadcmd(nickname, destination string) {
	err := callAPI(fmt.Sprintf("/renter/download?nickname=%s&destination=%s", nickname, destination))
	if err != nil {
		fmt.Println("Could not download file:", err)
		return
	}
	fmt.Printf("Started downloading '%s' to %s.\n", nickname, destination)
}

func renterdownloadqueuecmd() {
	var queue []api.DownloadInfo
	err := getAPI("/renter/downloadqueue", &queue)
	if err != nil {
		fmt.Println("Could not get download queue:", err)
		return
	}
	if len(queue) == 0 {
		fmt.Println("No downloads to show.")
		return
	}
	fmt.Println("Download Queue:")
	for _, file := range queue {
		fmt.Printf("%5.1f%% %s -> %s\n", 100*float32(file.Received)/float32(file.Filesize), file.Nickname, file.Destination)
	}
}

func renterstatuscmd() {
	status := new(modules.RentInfo)
	err := getAPI("/renter/status", status)
	if err != nil {
		fmt.Println("Could not get file status:", err)
		return
	}
	if len(status.Files) == 0 {
		fmt.Println("No files have been uploaded.")
		return
	}
	fmt.Println("Uploaded", len(status.Files), "files:")
	for _, file := range status.Files {
		fmt.Println("\t", file)
	}
}
