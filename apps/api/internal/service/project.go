package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/mokuhyo-driven-test/api/internal/model"
	"github.com/mokuhyo-driven-test/api/internal/repository"
)

type ProjectService struct {
	projectRepo *repository.ProjectRepository
	nodeRepo    *repository.NodeRepository
	edgeRepo    *repository.EdgeRepository
}

func NewProjectService(projectRepo *repository.ProjectRepository, nodeRepo *repository.NodeRepository, edgeRepo *repository.EdgeRepository) *ProjectService {
	return &ProjectService{
		projectRepo: projectRepo,
		nodeRepo:    nodeRepo,
		edgeRepo:    edgeRepo,
	}
}

func (s *ProjectService) CreateProject(ctx context.Context, userID uuid.UUID, req model.CreateProjectRequest) (*model.Project, error) {
	project, err := s.projectRepo.Create(ctx, userID, req)
	if err != nil {
		return nil, err
	}

	// Create initial root node
	initialNode, err := s.nodeRepo.Create(ctx, project.ID, req.Title)
	if err != nil {
		return nil, fmt.Errorf("failed to create initial node: %w", err)
	}

	// Create edge for root node
	_, err = s.edgeRepo.Create(ctx, project.ID, nil, initialNode.ID, model.RelationNeutral, nil, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to create initial edge: %w", err)
	}

	return project, nil
}

func (s *ProjectService) ListProjects(ctx context.Context, userID uuid.UUID) ([]model.Project, error) {
	return s.projectRepo.ListByUserID(ctx, userID)
}

func (s *ProjectService) GetProject(ctx context.Context, projectID uuid.UUID) (*model.Project, error) {
	return s.projectRepo.GetByID(ctx, projectID)
}

func (s *ProjectService) UpdateProject(ctx context.Context, projectID uuid.UUID, req model.UpdateProjectRequest) error {
	return s.projectRepo.Update(ctx, projectID, req)
}

func (s *ProjectService) CheckOwnership(ctx context.Context, projectID, userID uuid.UUID) (bool, error) {
	return s.projectRepo.CheckOwnership(ctx, projectID, userID)
}

func (s *ProjectService) SaveProject(ctx context.Context, projectID uuid.UUID) error {
	return s.projectRepo.UpdateUpdatedAt(ctx, projectID)
}

func (s *ProjectService) GetTree(ctx context.Context, projectID uuid.UUID) (*model.TreeResponse, error) {
	project, err := s.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if project == nil {
		return nil, fmt.Errorf("project not found")
	}

	nodes, err := s.nodeRepo.ListByProjectID(ctx, projectID)
	if err != nil {
		return nil, err
	}

	edges, err := s.edgeRepo.ListByProjectID(ctx, projectID)
	if err != nil {
		return nil, err
	}

	return &model.TreeResponse{
		Project: *project,
		Nodes:   nodes,
		Edges:   edges,
	}, nil
}
