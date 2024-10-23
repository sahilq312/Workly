package controller // Package declaration

import ( // Importing necessary packages
	"net/http" // For handling HTTP requests

	"github.com/gin-gonic/gin"                // For creating web applications
	"github.com/sahilq312/workly/initializer" // For initializing the application
	"github.com/sahilq312/workly/model"       // For defining the data model
)

func CreateJob(c *gin.Context) { // Function to create a new job
	// Get request body
	var body struct { // Define the structure of the request body
		Title       string   `json:"title"`       // Title of the job
		Description string   `json:"description"` // Description of the job
		Location    string   `json:"location"`    // Location of the job
		Salary      string   `json:"salary"`      // Salary of the job (as a string)
		CompanyID   uint     `json:"company_id"`  // ID of the company that posted the job
		Skills      []string `json:"skills"`      // Array of skills required for the job
	}

	// Bind JSON and check for errors
	if err := c.BindJSON(&body); err != nil { // Bind the request body to the defined structure
		c.JSON(http.StatusBadRequest, gin.H{ // If there's an error, return a bad request response
			"error": "Invalid request body", // Error message
		})
		return
	}

	// Validate required fields
	if body.Title == "" || body.Description == "" || body.Location == "" || body.Salary == "" || body.CompanyID == 0 || len(body.Skills) == 0 { // Check if any of the required fields are empty
		c.JSON(http.StatusBadRequest, gin.H{ // If any of the required fields are empty, return a bad request response
			"error": "All fields are required", // Error message
		})
		return
	}

	// Create the job entry
	job := model.Job{ // Create a new job instance
		Title:       body.Title,       // Set the title of the job
		Description: body.Description, // Set the description of the job
		Location:    body.Location,    // Set the location of the job
		Salary:      body.Salary,      // Set the salary of the job
		CompanyID:   body.CompanyID,   // Set the company ID of the job
		Skills:      body.Skills,      // Set the skills required for the job
	}
	result := initializer.DB.Create(&job) // Save the job to the database
	if result.Error != nil {              // If there's an error while saving the job
		c.JSON(http.StatusInternalServerError, gin.H{ // Return a server error response
			"error": "Failed to create job", // Error message
		})
		return
	}

	// Return the created job
	c.JSON(http.StatusOK, gin.H{ // If the job is created successfully, return a success response
		"message": "Job created successfully", // Success message
		"job":     job,                        // The created job
	})
}

func GetJob(c *gin.Context) { // Function to get a job
	// Get job
	id := c.Param("id") // Get the job ID from the request parameter
	// Return job
	c.JSON(http.StatusOK, gin.H{ // Return the job with the specified ID
		"message": id, // The job ID
	})
}

func UpdateJob(c *gin.Context) { // Function to update a job
	// Get job
	var body struct { // Define the structure of the request body
		Title       string   // Title of the job
		Description string   // Description of the job
		Skills      []string // Array of skills required for the job
		Location    string   // Location of the job
		Salary      string   // Salary of the job
		CompanyId   string   // ID of the company that posted the job
		jobId       string   // ID of the job to be updated
	}
	err := c.BindJSON(&body) // Bind the request body to the defined structure
	if err != nil {          // If there's an error while binding the request body
		c.JSON(http.StatusBadRequest, gin.H{ // Return a bad request response
			"error": "Invalid request body", // Error message
		})
		return
	}

	// Update job
	initializer.DB.Model(&model.Job{}).Where("id = ?", body.jobId).Updates(body) // Update the job with the specified ID
	// Return job
	c.JSON(http.StatusOK, gin.H{ // If the job is updated successfully, return a success response
		"message": "Job updated successfully", // Success message
	})
}

func DeleteJob(c *gin.Context) { // Function to delete a job
	// Get job
	var body struct { // Define the structure of the request body
		jobId string // ID of the job to be deleted
	}
	err := c.BindHeader(&body) // Bind the request header to the defined structure
	if err != nil {            // If there's an error while binding the request header
		c.JSON(http.StatusBadRequest, gin.H{ // Return a bad request response
			"error": "Invalid request body", // Error message
		})
		return
	}
	// Delete job
	initializer.DB.Delete(&model.Job{}, body.jobId) // Delete the job with the specified ID
	// Return job
	c.JSON(http.StatusOK, gin.H{ // If the job is deleted successfully, return a success response
		"message": "Job deleted successfully", // Success message
	})
}

func GetAllJobs(c *gin.Context) { // Function to get all jobs
	// Get all jobs
	jobs := []model.Job{}      // Get all jobs from the database
	initializer.DB.Find(&jobs) // Find all jobs in the database
	// Return all jobs
	c.JSON(http.StatusOK, gin.H{ // Return all the jobs
		"message": "Jobs fetched successfully", // Success message
		"jobs":    jobs,                        // The fetched jobs
	})
}

func GetJobsByCompany(c *gin.Context) { // Function to get jobs by company
	var body struct { // Define the structure of the request body
		companyId string // ID of the company to get jobs for
	}
	err := c.BindHeader(&body) // Bind the request header to the defined structure
	if err != nil {            // If there's an error while binding the request header
		c.JSON(http.StatusBadRequest, gin.H{ // Return a bad request response
			"error": "Invalid request body", // Error message
		})
		return
	}
	// Get jobs by company
	jobs := []model.Job{}                                              // Get all jobs from the database
	initializer.DB.Where("company_id = ?", body.companyId).Find(&jobs) // Find all jobs in the database for the specified company
	// Return jobs by company
	c.JSON(http.StatusOK, gin.H{ // Return all the jobs for the specified company
		"message": "Jobs fetched successfully", // Success message
		"jobs":    jobs,                        // The fetched jobs
	})
}
func GetJobsByLocation(c *gin.Context) { // Function to get jobs by location
	// Get jobs by location
	var body struct { // Define the structure of the request body
		Location string `json:"location"` // Location to get jobs for
	}
	err := c.BindJSON(&body) // Bind the request body to the defined structure
	if err != nil {          // If there's an error while binding the request body
		c.JSON(http.StatusBadRequest, gin.H{ // Return a bad request response
			"error": "Invalid request body", // Error message
		})
		return
	}

	// Return jobs by location
	jobs := []model.Job{}                                           // Get all jobs from the database
	initializer.DB.Where("location = ?", body.Location).Find(&jobs) // Find all jobs in the database for the specified location

	if len(jobs) == 0 { // If no jobs are found for the specified location
		c.JSON(http.StatusNotFound, gin.H{ // Return a not found response
			"error": "No jobs found", // Error message
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{ // Return all the jobs for the specified location
		"message": "Jobs fetched successfully", // Success message
		"jobs":    jobs,                        // The fetched jobs
	})
}

func GetJobsBySkill(c *gin.Context) { // Function to get jobs by skill
	// Get jobs by skill
	var body struct { // Define the structure of the request body
		Skill string `json:"skill"` // Skill to get jobs for
	}
	err := c.BindJSON(&body) // Bind the request body to the defined structure
	if err != nil {          // If there's an error while binding the request body
		c.JSON(http.StatusBadRequest, gin.H{ // Return a bad request response
			"error": "Invalid request body", // Error message
		})
		return
	}
	jobs := []model.Job{}                                                 // Get all jobs from the database
	initializer.DB.Where("skills LIKE ?", "%"+body.Skill+"%").Find(&jobs) // Find all jobs in the database that require the specified skill

	if len(jobs) == 0 { // If no jobs are found for the specified skill
		c.JSON(http.StatusNotFound, gin.H{ // Return a not found response
			"error": "No jobs found", // Error message
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{ // Return all the jobs for the specified skill
		"message": "Jobs fetched successfully", // Success message
		"jobs":    jobs,                        // The fetched jobs
	})
}
