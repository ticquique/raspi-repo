/*
 * File share api
 *
 * File share api.
 *
 * API version: 2.0.0
 * Contact: enponsba@gmail.com
 */

package raspi

import "os"

type InlineResponse200 struct {
	Total int64      `json:"total,omitempty"`
	Items []*os.File `json:"items,omitempty"`
}
