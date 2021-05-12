package main

import (
	"encoding/csv"
	"fmt"
	"os"
)

type csvData struct {
	path    string
	content [][]string
}

// getDataHTML returns the content of the csvData struct HTML table data
func (d *csvData) getDataHTML() string {
	var html string

	for i, row := range d.content {
		if i == 0 {
			html += "<thead>\n\t\t<tr>\n"
		} else {
			html += fmt.Sprintf("\t\t<tr>\n")
		}

		if i != 0 {
			for _, cell := range row {
				html += fmt.Sprintf("\t\t\t<td>%s</td>\n", cell)
			}
		} else {
			for _, cell := range row {
				html += fmt.Sprintf("\t\t\t<th>%s</th>\n", cell)
			}
		}

		if i != 0 {
			html += "\t\t</tr>\n"
		} else {
			html += "\t\t</tr>\n\t</thead>\n\t<tbody>\n"
		}
	}
	html += "\t</tbody>"
	return html
}

func (d *csvData) loadData() error {
	file, err := os.Open(d.path)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	d.content, err = reader.ReadAll()
	if err != nil {
		return err
	}

	return nil
}

func (d *csvData) add(cells []string) error {
	cells = append([]string{fmt.Sprintf("%d", len(d.content))}, cells...)
	d.content = append(d.content, cells)
	file, err := os.OpenFile(d.path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)

	err = writer.Write(cells)
	if err != nil {
		return err
	}
	defer writer.Flush()

	return nil
}

func (d *csvData) delete(index int) error {
	if index < 1 || index > len(d.content)-1 {
		return fmt.Errorf("index has to be bigger than 1 and can't be bigger than the last element")
	}

	d.content = append(d.content[:index], d.content[index+1:]...)

	file, err := os.Create(d.path)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)

	for _, rec := range d.content {
		err = writer.Write(rec)
		if err != nil {
			return err
		}
	}
	defer writer.Flush()

	return nil
}
