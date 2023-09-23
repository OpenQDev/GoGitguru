package util

import "os"

func DeleteLocalRepoAndTarball(prefixPath string, repo string) error {
	err := os.RemoveAll(prefixPath + repo)
	if err != nil {
		return err
	}
	err = os.RemoveAll(prefixPath + repo + ".tar.gz")
	if err != nil {
		return err
	}
	return nil
}
