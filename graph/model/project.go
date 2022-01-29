package model

import (
	"time"

	"gitlab.lrz.de/projecthub/gql-api/sql"
)

type Project struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Languages   []string   `json:"languages"`
	Location    *Location  `json:"location"`
	CreatedAt   *time.Time `json:"createdAt"`
	CreatorID   string
}

func ProjectFromDBProject(dbProject sql.Project) *Project {
	createdAt := &dbProject.CreatedAt.Time
	if !dbProject.CreatedAt.Valid {
		createdAt = nil
	}
	var location *Location
	if !dbProject.Location.Valid {
		location = &Location{
			Name: dbProject.Location.String,
		}
	}
	return &Project{
		ID:          dbProject.ID.String(),
		Name:        dbProject.Name,
		Description: dbProject.Description,
		CreatorID:   dbProject.Creator.String(),
		CreatedAt:   createdAt,
		Languages:   []string{},
		Location:    location,
	}
}

func ProjectsFromDBProjects(dbProject []sql.Project) []*Project {
	projects := make([]*Project, len(dbProject))
	for i, dbProject := range dbProject {
		projects[i] = ProjectFromDBProject(dbProject)
	}
	return projects
}

type ProjectMutation struct {
	ID string
}
