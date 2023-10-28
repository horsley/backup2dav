package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/studio-b12/gowebdav"
	"gopkg.in/yaml.v2"
)

var appConfig config

func main() {
	f, err := os.Open("config.yaml")
	if err != nil {
		log.Fatalln("open config file err:", err)
	}
	defer f.Close()
	rd := yaml.NewDecoder(f)
	err = rd.Decode(&appConfig)
	if err != nil {
		log.Fatalln("read config file err:", err)
	}

	var buf bytes.Buffer

	for _, job := range appConfig.ListJobs() {
		filename := fmt.Sprintf(`%s-%s.tar.gz`, job.Name, time.Now().Format(job.TimeFormat))

		err = createArchive(job.Dir, &buf)
		if err != nil {
			log.Fatalln("archive err:", err)
		}

		client := gowebdav.NewClient(job.WebDAV, job.User, job.Password)
		err = client.MkdirAll(job.Name, 0755)
		if err != nil {
			log.Fatalln("mkdir err:", err)
		}

		err = client.WriteStream(path.Join(job.Name, filename), &buf, 0644)
		if err != nil {
			log.Fatalln("upload err:", err)
		}
		log.Println("uploaded", filename)

		buf.Reset()

		err = rotateBackups(client, job.Name, job.Rotate, fmt.Sprintf(`%s-%s.tar.gz`, job.Name, job.TimeFormat))
		if err != nil {
			log.Fatalln("rotate err:", err)
		}
	}
}

func createArchive(dir string, buf *bytes.Buffer) error {
	gz := gzip.NewWriter(buf)
	defer gz.Close()

	tarball := tar.NewWriter(gz)
	defer tarball.Close()

	return filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		stat, err := d.Info()
		if err != nil {
			return err
		}

		header, err := tar.FileInfoHeader(stat, "")
		if err != nil {
			return err
		}
		header.Name = path

		err = tarball.WriteHeader(header)
		if err != nil {
			return err
		}

		_, err = io.Copy(tarball, file)

		log.Println("add to tarball:", path)
		return err
	})
}

func rotateBackups(client *gowebdav.Client, saveDir, duration, fileTimePattern string) error {
	if duration == "" {
		return nil
	}

	dur, err := time.ParseDuration(duration)
	if err != nil {
		return err
	}

	files, err := client.ReadDir(saveDir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if t, err := time.Parse(fileTimePattern, file.Name()); err != nil {
			log.Println("skip file not match:", file.Name())
			continue
		} else if !t.Before(time.Now().Add(-dur)) {
			continue
		}

		toRemove := filepath.Join(saveDir, file.Name())
		log.Println("remove file:", toRemove, "err:", client.Remove(toRemove))
	}

	return nil
}
