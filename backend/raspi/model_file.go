/*
 * File share api
 *
 * File share api.
 *
 * API version: 2.0.0
 * Contact: enponsba@gmail.com
 */

package raspi

type File struct {
	Id       int64  `json:"id" db:"id"`
	Title    string `json:"title" db:"title"`
	Filename string `json:"filename" db:"filename"`
	Route    string `json:"route" db:"route"`
	Type_    string `json:"type" db:"type"`
	Image    string `json:"image" db:"image"`
	Summary  string `json:"summary" db:"summary"`
}
