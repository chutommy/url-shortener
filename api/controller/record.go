package controller

import (
	"net/http"
	"strconv"

	"github.com/chutified/url-shortener/api/data"
	"github.com/gin-gonic/gin"
)

// AddRecord adds a new record.
func (h *handler) AddRecord(c *gin.Context) {

	// bind record
	var newr data.Record
	err := c.ShouldBindJSON(&newr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// add record
	r, err := h.ds.AddRecord(c, &newr)
	if err != nil {
		switch err {

		case data.ErrInvalidRecord:
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

		case data.ErrUnavailableShort:
			c.JSON(http.StatusConflict, gin.H{
				"error": err.Error(),
			})

		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}
		return
	}

	// record successfully added
	c.JSON(http.StatusOK, r)
}

// UpdateRecord replace the record with given the certain ID.
func (h *handler) UpdateRecord(c *gin.Context) {

	// get record's ID
	id := c.Param("record_id")

	// bind record
	var newr data.Record
	err := c.ShouldBindJSON(&newr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// update record
	r, err := h.ds.UpdateRecord(c, id, &newr)
	if err != nil {
		switch err {

		case data.ErrInvalidRecord:
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

		case data.ErrUnavailableShort:
			c.JSON(http.StatusConflict, gin.H{
				"error": err.Error(),
			})

		case data.ErrIDNotFound:
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})

		// server error
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}
		return
	}

	// record successfully updated
	c.JSON(http.StatusOK, r)
}

// DeleteRecord removes the record with the certain ID.
func (h *handler) DeleteRecord(c *gin.Context) {

	// get record's ID
	id := c.Param("record_id")

	// delete record
	did, err := h.ds.DeleteRecord(c, id)
	if err != nil {
		switch err {

		case data.ErrIDNotFound:
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})

		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}
		return
	}

	// record successfully deleted
	c.JSON(http.StatusOK, gin.H{
		"id": did,
	})
}

// GetRecordByID serves a record with the certain ID.
func (h *handler) GetRecordByID(c *gin.Context) {

	// get record's ID
	id := c.Param("record_id")

	// get record
	r, err := h.ds.GetRecordByID(c, id)
	if err != nil {
		switch err {

		case data.ErrIDNotFound:
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})

		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}
		return
	}

	// record successfully retrieved
	c.JSON(http.StatusOK, r)
}

// GetRecordByShort serves a record with the certain Short value.
func (h *handler) GetRecordByShort(c *gin.Context) {

	// get record's Short
	short := c.Param("record_short")

	// get record
	r, err := h.ds.GetRecordByShort(c, short)
	if err != nil {
		switch err {

		case data.ErrShortNotFound:
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})

		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}
		return
	}

	// record successfully retrieved
	c.JSON(http.StatusOK, r)
}

// GetRecordByFull serves a record with the sertain Full value.
func (h *handler) GetRecordByFull(c *gin.Context) {

	// get record's Full
	full := c.Param("record_full")

	// get record
	r, err := h.ds.GetRecordByFull(c, full)
	if err != nil {
		switch err {

		case data.ErrFullNotFound:
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})

		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}
		return
	}

	// record successfully retrieved
	c.JSON(http.StatusOK, r)
}

// GetRecordsLen returns a total number of records.
func (h *handler) GetRecordsLen(c *gin.Context) {

	// get length
	l, err := h.ds.GetRecordsLen(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	// records length successfully retrieved
	c.JSON(http.StatusOK, gin.H{
		"len": l,
	})
}

// GetAllRecords returns xth page with a certain number of records.
func (h *handler) GetAllRecords(c *gin.Context) {

	// get pagination data
	pageStr := c.DefaultQuery("page", "1")
	paginStr := c.DefaultQuery("pagin", "100")
	sort := c.DefaultQuery("sort", "created_at")

	// convert page and pagin
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 1
	}
	pagin, err := strconv.Atoi(paginStr)
	if err != nil {
		pagin = 30
	}

	// declare PageCfg
	pcfg := data.PageCfg{
		Page:  page,
		Pagin: pagin,
		Sort:  sort,
	}

	// get records
	rs, pgcfgG, err := h.ds.GetAllRecords(c, pcfg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	// records successfully retrieved
	c.JSON(http.StatusOK, gin.H{
		"page_config": pgcfgG,
		"result":      rs,
	})
}