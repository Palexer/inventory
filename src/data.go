package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
)

type csvData struct {
	contentPath string
	content     [][]string
	cachePath   string
	cache       [][]string
	key         string
	heading     string
}

// getDataHTML returns the content of the csvData struct HTML table data
func (d *csvData) getDataHTML() string {
	if err := d.loadData(); err != nil {
		log.Printf("failed to convert CSV to HTML data: %v\n", err)
	}
	var html string

	for i, row := range d.content {
		if i == 0 {
			html += "<thead>\n\t\t\t<tr>\n\t\t\t\t<th>Nr.</th>\n"

			for _, cell := range row {
				html += fmt.Sprintf("\t\t\t\t<th>%s</th>\n", cell)
			}
			html += "\t\t\t\t<td></td>\n\t\t\t\t<td></td>\n\t\t\t</tr>\n\t\t</thead>\n\t\t<tbody>\n"

		} else {
			html += fmt.Sprintf("\t\t\t<tr>\n\t\t\t\t<td>%d</td>\n", i)

			for _, cell := range row {
				html += fmt.Sprintf("\t\t\t\t<td>%s</td>\n", cell)
			}
			html += "\t\t\t\t<td class=\"editCell\"><i class=\"far fa-edit\"></i></td>\n\t\t\t\t<td class=\"deleteCell\"><i class=\"far fa-trash-alt\"></i></td>\n\t\t\t</tr>\n"
		}
	}
	html += "\t</tbody>"
	return html
}

func (d *csvData) loadData() error {
	// create data and cache file if they don't exist
	if _, err := os.Stat(data.contentPath); os.IsNotExist(err) {
		_, err := os.Create(data.contentPath)
		if err != nil {
			return fmt.Errorf("failed to create data file: %v\n", err)
		}
		// default table heading
		data.add([]string{"Name", "Description", "Count", "Date"})
	}

	if _, err := os.Stat(data.cachePath); os.IsNotExist(err) {
		_, err := os.Create(data.cachePath)
		if err != nil {
			return fmt.Errorf("failed to create cache file: %v\n", err)
		}
	}

	fileContent, err := os.Open(d.contentPath)
	if err != nil {
		return err
	}
	defer fileContent.Close()

	reader := csv.NewReader(fileContent)
	d.content, err = reader.ReadAll()
	if err != nil {
		return err
	}

	fileCache, err := os.Open(d.cachePath)
	if err != nil {
		return err
	}
	defer fileCache.Close()

	reader = csv.NewReader(fileCache)
	d.cache, err = reader.ReadAll()
	if err != nil {
		return err
	}

	return nil
}

func (d *csvData) add(cells []string) error {
	if err := d.loadData(); err != nil {
		return err
	}
	d.content = append(d.content, cells)

	// write file
	file, err := os.OpenFile(d.contentPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
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

func (d *csvData) edit(index int, newCells []string) error {
	if err := d.loadData(); err != nil {
		return err
	}
	d.content[index] = newCells

	// write file
	file, err := os.Create(d.contentPath)
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

func (d *csvData) delete(index int) error {
	if err := d.loadData(); err != nil {
		return err
	}
	if index < 1 || index > len(d.content)-1 {
		return fmt.Errorf("index has to be bigger than 1 and can't be bigger than the last element")
	}

	d.cache = append(d.cache, d.content[index])
	d.content = append(d.content[:index], d.content[index+1:]...)

	// write file
	fileContent, err := os.Create(d.contentPath)
	if err != nil {
		return err
	}
	defer fileContent.Close()

	writer := csv.NewWriter(fileContent)

	for _, rec := range d.content {
		err = writer.Write(rec)
		if err != nil {
			return err
		}
	}
	defer writer.Flush()

	fileCache, err := os.Create(d.cachePath)
	if err != nil {
		return err
	}
	defer fileCache.Close()

	writer = csv.NewWriter(fileCache)

	for _, rec := range d.cache {
		err = writer.Write(rec)
		if err != nil {
			return err
		}
	}
	defer writer.Flush()

	return nil
}

func (d *csvData) restore() error {
	if err := d.loadData(); err != nil {
		return err
	}
	if len(d.cache) < 1 {
		return fmt.Errorf("failed to restore: empty cache")
	}

	cells := d.cache[len(d.cache)-1]
	d.content = append(d.content, cells)
	d.cache = d.cache[:len(d.cache)-1]

	// write content file
	file, err := os.OpenFile(d.contentPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
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

	// write cache file
	fileCache, err := os.Create(d.cachePath)
	if err != nil {
		return err
	}
	defer fileCache.Close()

	writer = csv.NewWriter(fileCache)

	for _, rec := range d.cache {
		err = writer.Write(rec)
		if err != nil {
			return err
		}
	}
	defer writer.Flush()
	return nil
}

// deleteAllAndBackUp deletes the currently used CSV file and creates a backup of it
func (d *csvData) deleteAllAndBackUp() error {
	if err := d.loadData(); err != nil {
		return err
	}

	err := os.Remove(d.contentPath)
	if err != nil {
		return err
	}

	if _, err := os.Stat("inventory_backup.csv"); err == nil {
		err := os.Remove("inventory_backup.csv")
		if err != nil {
			return err
		}
	}

	file, err := os.Create("inventory_backup.csv")
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
