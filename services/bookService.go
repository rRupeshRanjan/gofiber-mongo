package services

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"gofiber-mongo/config"
	"gofiber-mongo/domain"
	"gofiber-mongo/repository"
	"strconv"
)

var logger = config.AppLogger

func GetBookByIdHandler(c *fiber.Ctx) error {
	c.Set("Content-Type", "application/json")

	id, requestResponse, validId := extractIdFromParamsIfValid(c)
	if validId {
		book, getBookResponse, bookFound := getBookByIdIfPresent(c, id)
		requestResponse = getBookResponse
		if bookFound {
			bookString, marshalError := json.Marshal(book)
			if marshalError != nil {
				logger.Error("Error converting book data to json: " + marshalError.Error())
				return c.Status(500).Send(nil)
			}
			requestResponse = c.Status(200).Send(bookString)
		}
	}
	return requestResponse
}

func CreateBookHandler(c *fiber.Ctx) error {
	c.Set("Content-Type", "application/json")

	book, requestResponse, validBook := getBookIfValid(c)
	if validBook {
		bookId, createBookError := repository.CreateBook(book)
		if createBookError != nil {
			logger.Error("Error while creating book: " + createBookError.Error())
			return c.Status(500).Send(nil)
		}

		book.Id = bookId
		bookString, unmarshalError := json.Marshal(book)
		if unmarshalError != nil {
			logger.Error("Error converting data to json format: " + unmarshalError.Error())
			return c.Status(500).Send(nil)
		}
		requestResponse = c.Status(200).Send(bookString)
	}
	return requestResponse
}

func UpdateBookHandler(c *fiber.Ctx) error {
	c.Set("Content-Type", "application/json")

	id, requestResponse, validId := extractIdFromParamsIfValid(c)
	if validId {
		book, requestValidityCheckResponse, validRequest := isValidUpdateRequest(c, id)
		requestResponse = requestValidityCheckResponse

		if validRequest {
			_, bookPresenceCheckResponse, bookFound := getBookByIdIfPresent(c, id)
			requestResponse = bookPresenceCheckResponse

			if bookFound {
				updateBookError := repository.UpdateBook(book)
				if updateBookError != nil {
					logger.Error("Error while updating book: " + updateBookError.Error())
					return c.Status(500).Send(nil)
				}
				requestResponse = c.Status(200).Send(c.Body())
			}
		}
	}
	return requestResponse
}

func DeleteBookHandler(c *fiber.Ctx) error {
	id, requestResponse, valid := extractIdFromParamsIfValid(c)
	if valid {
		deleteBookError := repository.DeleteBookById(id)
		if deleteBookError != nil {
			logger.Error("Error deleting book with id: " + strconv.FormatInt(id, 10))
			return c.Status(500).Send(nil)
		}
		requestResponse = c.Status(204).Send(nil)
	}
	return requestResponse
}

func getBookByIdIfPresent(c *fiber.Ctx, id int64) (domain.Book, error, bool) {
	book, err := repository.GetBookById(id)
	if err != nil {
		if err.Error() == domain.NoDocs {
			return domain.Book{}, c.Status(404).Send(nil), false
		}
		fmt.Println(err.Error())
		return domain.Book{}, c.Status(500).Send(nil), false
	}
	return book, nil, true
}

func extractIdFromParamsIfValid(c *fiber.Ctx) (int64, error, bool) {
	id, err := strconv.ParseInt(c.Params("id"), 10, 64)
	if err != nil {
		logger.Error("Error fetching integer from id: " + c.Params("id"))
		return -1, c.Status(400).Send(nil), false
	}
	return id, nil, true
}

func getBookIfValid(c *fiber.Ctx) (domain.Book, error, bool) {
	var book domain.Book
	unmarshalError := json.Unmarshal(c.Body(), &book)
	if unmarshalError != nil {
		logger.Error("Error converting body to book: " + unmarshalError.Error())
		return domain.Book{}, c.Status(400).Send(nil), false
	}
	return book, nil, true
}

func isValidUpdateRequest(c *fiber.Ctx, id int64) (domain.Book, error, bool) {
	book, response, valid := getBookIfValid(c)
	if !valid {
		return domain.Book{}, response, false
	} else if book.Id != id {
		logger.Error("Book id from URL: " + strconv.FormatInt(id, 10) + " is different from id in request body: " + strconv.FormatInt(book.Id, 10))
		return domain.Book{}, c.Status(400).Send(nil), false
	}
	return book, nil, true
}
