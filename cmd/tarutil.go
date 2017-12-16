package main

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// TarWorkdir tars and gzips a directory and returns the absolute path
// of the generated archive as a string or, in case, an error
func TarWorkdir(source string) (string, error) {
	logger.Debugf("packing directory: %s", source)
	tarfile, err := ioutil.TempFile("", "hydra-")
	defer tarfile.Close()
	check(err)
	logger.Infof("created archive: %s", tarfile.Name())

	gzipWriter := gzip.NewWriter(tarfile)
	defer gzipWriter.Close()

	tarball := tar.NewWriter(gzipWriter)
	defer tarball.Close()

	return tarfile.Name(), filepath.Walk(source,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// You get an error if you have just a directory
			if info.IsDir() {
				return nil
			}

			header, err := tar.FileInfoHeader(info, info.Name())
			if err != nil {
				return err
			}

			header.Name = filepath.Join(strings.TrimPrefix(path, source))
			logger.Debugf("adding file to archive: %s", header.Name)

			if err := tarball.WriteHeader(header); err != nil {
				return err
			}

			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = io.Copy(tarball, file)
			return err
		},
	)
}
