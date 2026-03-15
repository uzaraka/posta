/*
 *  MIT License
 *
 * Copyright (c) 2026 Jonas Kaninda
 *
 *  Permission is hereby granted, free of charge, to any person obtaining a copy
 *  of this software and associated documentation files (the "Software"), to deal
 *  in the Software without restriction, including without limitation the rights
 *  to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 *  copies of the Software, and to permit persons to whom the Software is
 *  furnished to do so, subject to the following conditions:
 *
 *  The above copyright notice and this permission notice shall be included in all
 *  copies or substantial portions of the Software.
 *
 *  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 *  IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 *  FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 *  AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 *  LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 *  OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 *  SOFTWARE.
 */

package handlers

import (
	"net/http"

	"github.com/jkaninda/okapi"
	"github.com/jkaninda/posta/internal/dto"
)

func ok[T any](c *okapi.Context, data T) error {
	return c.JSON(http.StatusOK, dto.Response[T]{
		Success: true,
		Data:    data,
	})
}

func created[T any](c *okapi.Context, data T) error {
	return c.JSON(http.StatusCreated, dto.Response[T]{
		Success: true,
		Data:    data,
	})
}

func noContent(c *okapi.Context) error {
	return c.JSON(http.StatusNoContent, dto.Response[any]{
		Success: true,
	})
}

func paginated[T any](c *okapi.Context, items []T, total int64, page, size int) error {
	totalPages := 0
	if size > 0 {
		totalPages = int((total + int64(size) - 1) / int64(size))
	}
	return c.JSON(http.StatusOK, dto.PageableResponse[T]{
		Success: true,
		Data:    items,
		Pageable: dto.Pageable{
			CurrentPage:   page,
			Size:          size,
			TotalPages:    totalPages,
			TotalElements: total,
			Empty:         len(items) == 0,
		},
	})
}
