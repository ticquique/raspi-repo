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
	Id    int    `json:"Id" db:"Id"`
	Alias string `json:"Alias" db:"Alias"`
	Name  string `json:"Name" db:"Name"`
	Route string `json:"Route" db:"Route"`
	Type_ string `json:"Type" db:"Type"`
}
