package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/mokuhyo-driven-test/api/internal/model"
	"github.com/mokuhyo-driven-test/api/internal/repository"
)

type NodeService struct {
	nodeRepo *repository.NodeRepository
	edgeRepo *repository.EdgeRepository
}

func NewNodeService(nodeRepo *repository.NodeRepository, edgeRepo *repository.EdgeRepository) *NodeService {
	return &NodeService{
		nodeRepo: nodeRepo,
		edgeRepo: edgeRepo,
	}
}

func (s *NodeService) CreateNode(ctx context.Context, projectID uuid.UUID, req model.CreateNodeRequest) (*model.Node, *model.Edge, error) {
	// Create node
	node, err := s.nodeRepo.Create(ctx, projectID, req.Content)
	if err != nil {
		return nil, nil, err
	}

	// Determine order_index
	orderIndex := 0
	if req.OrderIndex != nil {
		orderIndex = *req.OrderIndex
	} else {
		maxOrder, err := s.nodeRepo.GetMaxOrderIndex(ctx, projectID, req.ParentNodeID)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to get max order index: %w", err)
		}
		orderIndex = maxOrder
	}

	// Determine relation
	relation := model.RelationNeutral
	if req.Relation != "" {
		relation = model.RelationType(req.Relation)
	}

	// Create edge
	edge, err := s.edgeRepo.Create(ctx, projectID, req.ParentNodeID, node.ID, relation, req.RelationLabel, orderIndex)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create edge: %w", err)
	}

	return node, edge, nil
}

func (s *NodeService) UpdateNode(ctx context.Context, nodeID uuid.UUID, req model.UpdateNodeRequest) error {
	return s.nodeRepo.Update(ctx, nodeID, req.Content)
}

func (s *NodeService) DeleteNode(ctx context.Context, projectID, nodeID uuid.UUID) error {
	return s.nodeRepo.SoftDeleteWithDescendants(ctx, projectID, nodeID)
}
