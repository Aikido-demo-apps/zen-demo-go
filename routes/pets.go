package routes

import (
	"net/http"

	"zen-demo-go/database"

	"github.com/gin-gonic/gin"
)

// SetupPetRoutes configures pet API routes
func SetupPetRoutes(r *gin.Engine) {
	r.GET("/api/pets/", getAllPets)
	r.GET("/api/pets/:id", getPetByID)
	r.POST("/api/create", createPet)
	r.GET("/clear", clearPets)
}

func getAllPets(c *gin.Context) {
	pets := database.GetAllPets()
	c.JSON(http.StatusOK, pets)
}

func getPetByID(c *gin.Context) {
	id := c.Param("id")
	pet := database.GetPetByID(id)
	if pet == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pet not found"})
		return
	}
	c.JSON(http.StatusOK, pet)
}

func createPet(c *gin.Context) {
	var req struct {
		Name string `json:"name"`
	}
	if err := c.BindJSON(&req); err != nil {
		c.String(http.StatusBadRequest, "Invalid request")
		return
	}

	rowsCreated := database.CreatePetByName(req.Name)
	if rowsCreated == -1 {
		c.String(http.StatusInternalServerError, "Database error occurred")
		return
	}

	c.String(http.StatusOK, "Success!")
}

func clearPets(c *gin.Context) {
	database.ClearAll()
	c.String(http.StatusOK, "Cleared successfully.")
}
