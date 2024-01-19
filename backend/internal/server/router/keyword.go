package router

import (
	"encoding/csv"
	"net/http"
	"scraping-keyword-web/backend/pkg/db"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func (r *router) getKeywords(c *gin.Context) {
	ctx := c.Request.Context()
	size, _ := strconv.Atoi(c.Query("size"))
	page, _ := strconv.Atoi(c.Query("page"))
	sort := c.Query("sort")

	limit := size
	offset := (page - 1) * size

	claims, err := verifyTokenFromRequest(c.Request)
	if err != nil {
		log.Error().Msgf("failed to verify token due to %v", err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	userKeywords, total, err := r.userkeywordRepo.GetUserKeywordByUserID(ctx, claims.UserID, limit, offset, sort)
	if err != nil {
		log.Error().Msgf("failed to collect user keywords due to %v", err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	keywordIDs := make([]string, 0, len(userKeywords))
	for _, userKeyword := range userKeywords {
		keywordIDs = append(keywordIDs, userKeyword.KeywordID)
	}

	keywords, err := r.keywordRepo.GetKeywordResultByIDs(ctx, keywordIDs)
	if err != nil {
		log.Error().Msgf("failed to collect keyword results due to %v", err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, GetKeywordsResponse{
		KeywordResults: keywords,
		Metadata:       GetKeywordsResponseMetaData{Total: total},
	})
}

func (r *router) uploadKeywordFile(c *gin.Context) {
	ctx := c.Request.Context()
	file, _ := c.FormFile("file")

	claims, err := verifyTokenFromRequest(c.Request)
	if err != nil {
		log.Error().Msgf("failed to verify token due to %v", err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// Open the uploaded file
	uploadedFile, err := file.Open()
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to open the uploaded file")
		return
	}
	defer uploadedFile.Close()

	// Convert file content to string slice
	keywords := []string{}
	reader := csv.NewReader(uploadedFile)
	for {
		record, err := reader.Read()
		if err != nil {
			break
		}

		if len(record) > 0 {
			kw := strings.TrimSpace(record[0])
			keywords = append(keywords, kw)
		}
	}

	mapKeywordResults := make(map[string]db.KeywordResult, len(keywords))
	existedKeywordResult, _ := r.keywordRepo.GetKeywordResultByKeywords(ctx, keywords)
	if len(existedKeywordResult) > 0 {
		for _, kw := range existedKeywordResult {
			mapKeywordResults[kw.Keyword] = kw
		}
	}

	notExistedKeywords := make([]string, 0, len(keywords))
	for _, k := range keywords {
		if _, ok := mapKeywordResults[k]; !ok {
			notExistedKeywords = append(notExistedKeywords, k)
		}
	}

	if len(notExistedKeywords) == 0 {
		c.JSON(http.StatusOK, existedKeywordResult)
		return
	}

	// Create new keyword results to DB
	notExistedKeywordInfos, err := r.keywordRepo.CrawlKeywordResults(notExistedKeywords)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	newKeywordResults, err := r.keywordRepo.CreateKeywordResults(ctx, notExistedKeywordInfos)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	newKeywordIDs := make([]string, 0, len(newKeywordResults))
	for _, k := range newKeywordResults {
		newKeywordIDs = append(newKeywordIDs, k.KeywordID)
	}

	// Create new user-keyword relation data
	_, err = r.userkeywordRepo.CreateMultipleUserKeywordByUserID(ctx, claims.UserID, newKeywordIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, newKeywordResults)
}
