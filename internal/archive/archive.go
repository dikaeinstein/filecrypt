package archive

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

var ErrFailedToCompress = errors.New("filecrypt: failed to compress")

// PackFiles walks the specified paths and archives them; the resulting
// archive is gzipped with the maximum compression.
func PackFiles(paths []string) ([]byte, error) {
	var buf bytes.Buffer

	zw, err := gzip.NewWriterLevel(&buf, gzip.BestCompression)
	if err != nil {
		return nil, err
	}

	tw := tar.NewWriter(zw)

	for _, p := range paths {
		if err := walkPath(p, tw); err != nil {
			return nil, err
		}
	}

	tw.Close()
	zw.Close()

	return buf.Bytes(), nil
}

func walkPath(p string, tw *tar.Writer) error {
	return filepath.Walk(p, func(path string, info os.FileInfo, err error) error {
		// prevent panic by handling failure accessing a path
		if err != nil {
			fmt.Println("err reading path")
			return err
		}

		if !info.Mode().IsDir() && !info.Mode().IsRegular() {
			return ErrFailedToCompress
		}

		fp := filepath.Clean(path)

		hdr, err := tar.FileInfoHeader(info, fp)
		if err != nil {
			return err
		}

		hdr.Name = fp
		if err = tw.WriteHeader(hdr); err != nil {
			return err
		}

		if info.Mode().IsRegular() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			_, err = io.Copy(tw, file)
			return err
		}

		return nil
	})
}

// UnpackFiles decompresses and unarchives the gzipped tarball passed in
// as data and uncompresses it to the top-level directory given. If
// unpack is false, only file names will be listed.
func UnpackFiles(in []byte, top string, unpack bool) error {
	buf := bytes.NewBuffer(in)

	zr, err := gzip.NewReader(buf)
	if err != nil {
		return err
	}
	defer zr.Close()

	tr := tar.NewReader(zr)

	for {
		hdr, err := tr.Next()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return err
		}

		if !unpack {
			fmt.Println(hdr.Name)
			continue
		}

		fp := filepath.Clean(filepath.Join(top, hdr.Name))

		switch hdr.Typeflag {
		case tar.TypeReg:
			file, err := os.Create(fp)
			if err != nil {
				return err
			}

			_, err = io.Copy(file, tr)
			if err != nil {
				return err
			}

			if err := file.Chmod(fs.FileMode(hdr.Mode)); err != nil {
				return err
			}

			file.Close()
		case tar.TypeDir:
			if err := os.MkdirAll(fp, fs.FileMode(hdr.Mode)); err != nil {
				return err
			}
		}
	}
}
