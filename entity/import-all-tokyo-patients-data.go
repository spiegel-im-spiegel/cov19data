package entity

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"io"

	"github.com/spiegel-im-spiegel/cov19data/ecode"
	"github.com/spiegel-im-spiegel/errs"
)

//ImportTokyoCSV function creates list of TokyoData from CSV.
func ImportTokyoCSV(r io.Reader, opts ...FiltersOptFunc) ([]TokyoData, error) {
	filter := NewFilters(opts...)
	records := []TokyoData{}
	cr := NewCsvReaderTokyo(r)
	for {
		record, err := cr.Next()
		if err != nil {
			if errs.Is(err, ecode.ErrNoData) {
				break
			}
			return nil, errs.Wrap(err)
		}
		if record.CheckFilter(filter) {
			records = append(records, record)
		}
	}
	return records, nil
}

func ExportTokyoJSON(data []TokyoData) ([]byte, error) {
	if len(data) == 0 {
		return nil, errs.Wrap(ecode.ErrNoData)
	}
	return json.Marshal(data)
}

//ExportWHOCSV function returns CSV string from list of WHOGlobalData.
func ExportTokyoCSV(data []TokyoData) ([]byte, error) {
	if len(data) == 0 {
		return nil, errs.Wrap(ecode.ErrNoData)
	}
	buf := &bytes.Buffer{}
	cw := csv.NewWriter(buf)
	cw.Comma = ','
	if err := cw.Write([]string{
		"発生日付",
		"地方公共団体コード",
		"対象者の居住地",
		"対象者の年代",
		"対象者の性別",
		"退院フラグ",
	}); err != nil {
		return nil, errs.Wrap(err)
	}
	for _, d := range data {
		if err := cw.Write([]string{
			d.Date.String(),
			d.LocalGovCode,
			d.Address,
			d.Age,
			d.Gender,
			d.LeaveFlag,
		}); err != nil {
			return nil, errs.Wrap(err)
		}
	}
	cw.Flush()
	return buf.Bytes(), nil
}

/* Copyright 2020 Spiegel
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * 	http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
