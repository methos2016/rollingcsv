package rollingcsv

import (
  "os"
  "encoding/csv"
  "strconv"
)

type RollingCsv struct {
  dir                string
  title              string
  headers          []string
  writeHeadersAll    bool
  maxBytes           int64
  maxLines           int64
  currentFile       *os.File
  currentFileWriter *csv.Writer
  currentFileNum     int64
  currentFileLines   int64
  currentFileBytes   int64
  files            []string
}

func New(title string, dir string, maxBytes int64, maxLines int64) (rCsv *RollingCsv) {
  rCsv = new(RollingCsv)
  rCsv.dir = dir
  rCsv.title = title
  rCsv.maxBytes = maxBytes
  rCsv.maxLines = maxLines
  rCsv.currentFile = nil
  rCsv.currentFileWriter = nil

  return rCsv
}


func (rCsv *RollingCsv) SetHeaders(headers []string, writeToAllFiles bool) {
  rCsv.headers = headers
  rCsv.writeHeadersAll = writeToAllFiles
}

func (rCsv *RollingCsv) Write(record []string) (err error) {
  var bytes int64 = 0
  for i := 0; i < len(record); i++ {
    bytes += int64(len(record[i]))
  }

  // we don't subtract one to compensate for the newline byte
  bytes += int64(len(record))

  if (rCsv.maxLines != 0 && rCsv.currentFileLines + 1 > rCsv.maxLines) ||
     (rCsv.maxBytes != 0 && rCsv.currentFileBytes + bytes > rCsv.maxBytes) ||
     rCsv.currentFile == nil {

    err = rCsv.NextFile()
    if err != nil {
      return err
    }
  }

  err = rCsv.currentFileWriter.Write(record)
  if err != nil {
    return err
  }

  rCsv.currentFileBytes += bytes
  rCsv.currentFileLines += 1

  return nil
}


func (rCsv *RollingCsv) GetOutputFiles() (files []string) {
  return rCsv.files
}

func (rCsv *RollingCsv) GetCurrentFileNumber() (n int64) {
  return rCsv.currentFileNum - 1
}


func (rCsv *RollingCsv) GetNextFileNumber() (n int64) {
  return rCsv.currentFileNum
}


func (rCsv *RollingCsv) NextFile() (err error) {
  err = rCsv.currentFile.Close()
  if err != nil {
    return err
  }

  fileNumStr := strconv.FormatInt(rCsv.currentFileNum, 10)
  fileName := rCsv.dir + "/" + rCsv.title + "-" + fileNumStr + ".csv"
  rCsv.currentFile, err = os.Create(fileName)
  if err != nil {
    return err
  }

  rCsv.currentFileWriter = csv.NewWriter(rCsv.currentFile)
  if len(rCsv.headers) != 0 && rCsv.writeHeadersAll {
    rCsv.Write(rCsv.headers)
  }

  rCsv.files = append(rCsv.files, fileName)
  rCsv.currentFileNum++
  rCsv.currentFileLines = 0
  rCsv.currentFileBytes = 0

  return nil
}

func (rCsv *RollingCsv) Close() (err error) {
  return rCsv.currentFile.Close()
}
