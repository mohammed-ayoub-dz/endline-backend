package handlers

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/mohammed-ayoub-js/endline-backend/internal/models"
	"github.com/mohammed-ayoub-js/endline-backend/internal/search"
	"gorm.io/gorm"
)

var (
    videoIndex *search.Index
    once       sync.Once
)

type SaveVideoRequest struct {
    Title       string `json:"title"`
    Description string `json:"description"`
    URL         string `json:"url"`
    Thumbnail   string `json:"thumbnail"`
    SubjectID   uint   `json:"subjectId"`
}

func LoadVideos() ([]models.Video, error) {
	filePath := "videos.json"
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		filePath = "internal/handler/videos.json" 
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("فشل قراءة الملف من المسار %s: %v", filePath, err)
	}

	var videos []models.Video
	if err := json.Unmarshal(data, &videos); err != nil {
		return nil, err
	}
	return videos, nil
}

func buildSearchIndex() {
    videos, err := LoadVideos() 
    if err != nil {
        fmt.Errorf("failed to load videos for indexing: %w", err)
    }
    idx := search.NewIndex()
    for _, v := range videos {
        idx.AddDocument(search.Video{
            ID:          int(v.ID),
            Title:       v.Title,
            Description: v.Description,
            URL:         v.URL,
            Thumbnail:   v.Thumbnail,
        })
    }
    idx.Build()
    videoIndex = idx
}

func SearchVideos(c *fiber.Ctx) error {
    once.Do(buildSearchIndex) 

    query := strings.TrimSpace(c.Query("q"))
    if query == "" {
        return c.Status(400).JSON(fiber.Map{"error": "Missing query"})
    }

    results := videoIndex.Search(query, 50)
    resp := make([]models.Video, len(results))
    for i, r := range results {
        resp[i] = models.Video{
            ID:          uint(r.Video.ID),
            Title:       r.Video.Title,
            Description: r.Video.Description,
            URL:         r.Video.URL,
            Thumbnail:   r.Video.Thumbnail,
        }
    }

    return c.JSON(fiber.Map{
        "query":   query,
        "count":   len(resp),
        "results": resp,
    })
}



func SaveVideoToDB(db *gorm.DB) fiber.Handler {
    return func(c *fiber.Ctx) error {
        var req SaveVideoRequest
        if err := c.BodyParser(&req); err != nil {
            return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
        }

         if req.SubjectID == 0 {
            var defaultSubject models.Subject
            err := db.Where("title = ?", "غير مصنف").First(&defaultSubject).Error
            if err != nil {
                defaultSubject = models.Subject{Title: "غير مصنف"}
                db.Create(&defaultSubject)
            }
            req.SubjectID = defaultSubject.ID
        }

        var existing models.Video
        err := db.Where("url = ?", req.URL).First(&existing).Error
        if err == nil {
            return c.JSON(fiber.Map{
                "message": "Video already exists",
                "videoId": existing.ID,
                "url":     existing.URL,
            })
        }


        video := models.Video{
            Title:       req.Title,
            Description: req.Description,
            URL:         req.URL,
            Thumbnail:   req.Thumbnail,
            SubjectID:   req.SubjectID,
        }
        if err := db.Create(&video).Error; err != nil {
            return c.Status(500).JSON(fiber.Map{"error": "Could not save video"})
        }

        return c.Status(201).JSON(fiber.Map{
            "message": "Video saved to library",
            "videoId": video.ID,
            "url":     video.URL,
        })
    }
}


func GetLibraryVideos(db *gorm.DB) fiber.Handler {
    return func(c *fiber.Ctx) error {
        var videos []models.Video
        if err := db.Order("created_at desc").Find(&videos).Error; err != nil {
            return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch library"})
        }

        fmt.Println(videos)
        return c.JSON(fiber.Map{
            "videos": videos,
            "count":  len(videos),
        })
    }
}