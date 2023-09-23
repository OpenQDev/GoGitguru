package main

func main() {
	organization := "OpenQDev"
	repo := "OpenQ-Workflows"
	cloneRepoAndUploadTarballToS3(organization, repo)
}
