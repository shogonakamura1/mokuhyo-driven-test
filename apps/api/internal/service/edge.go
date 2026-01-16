package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/mokuhyo-driven-test/api/internal/model"
	"github.com/mokuhyo-driven-test/api/internal/repository"
)

type EdgeService struct {
	edgeRepo *repository.EdgeRepository
}

func NewEdgeService(edgeRepo *repository.EdgeRepository) *EdgeService {
	return &EdgeService{edgeRepo: edgeRepo}
}

func (s *EdgeService) UpdateEdge(ctx context.Context, edgeID uuid.UUID, req model.UpdateEdgeRequest) error {
	return s.edgeRepo.Update(ctx, edgeID, req.Relation, req.RelationLabel)
}

func (s *EdgeService) Reorder(ctx context.Context, projectID uuid.UUID, req model.ReorderRequest) error {
	return s.edgeRepo.Reorder(ctx, projectID, req.ParentNodeID, req.OrderedChildNodeIDs)
}
