package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/mokuhyo-driven-test/api/internal/ai"
	"github.com/mokuhyo-driven-test/api/internal/model"
	"github.com/mokuhyo-driven-test/api/internal/repository"
)

type NodeService struct {
	nodeRepo         repository.NodeRepository
	edgeRepo         repository.EdgeRepository
	questionSelector ai.QuestionSelector
}

func NewNodeService(nodeRepo repository.NodeRepository, edgeRepo repository.EdgeRepository, questionSelector ai.QuestionSelector) *NodeService {
	return &NodeService{
		nodeRepo:         nodeRepo,
		edgeRepo:         edgeRepo,
		questionSelector: questionSelector,
	}
}

func (s *NodeService) CreateNode(ctx context.Context, projectID uuid.UUID, req model.CreateNodeRequest) (*model.Node, *model.Edge, error) {
	var question *string
	if req.ParentNodeID != nil {
		selected, err := s.selectQuestion(ctx, projectID, *req.ParentNodeID)
		if err != nil {
			fallback := fallbackQuestionSelection(defaultQuestionCandidates(), nil, nil)
			selected = fallback
		}
		question = &selected
	}

	// Create node
	node, err := s.nodeRepo.Create(ctx, projectID, req.Content, question)
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

func (s *NodeService) selectQuestion(ctx context.Context, projectID uuid.UUID, parentNodeID uuid.UUID) (string, error) {
	nodes, err := s.nodeRepo.ListByProjectID(ctx, projectID)
	if err != nil {
		return "", fmt.Errorf("failed to list nodes: %w", err)
	}
	edges, err := s.edgeRepo.ListByProjectID(ctx, projectID)
	if err != nil {
		return "", fmt.Errorf("failed to list edges: %w", err)
	}

	nodeByID := make(map[uuid.UUID]model.Node, len(nodes))
	for _, node := range nodes {
		nodeByID[node.ID] = node
	}
	parentNode, ok := nodeByID[parentNodeID]
	if !ok {
		return "", fmt.Errorf("parent node not found")
	}

	parentByChild := make(map[uuid.UUID]*uuid.UUID, len(edges))
	for _, edge := range edges {
		parentByChild[edge.ChildNodeID] = edge.ParentNodeID
	}

	ancestors := collectAncestors(parentNodeID, parentByChild, nodeByID)
	siblings := collectSiblings(parentNodeID, edges, nodeByID)

	candidates := questionCandidates()
	edgeByChild := make(map[uuid.UUID]model.Edge, len(edges))
	for _, edge := range edges {
		edgeByChild[edge.ChildNodeID] = edge
	}
	prompt := buildQuestionPrompt(parentNode, ancestors, siblings, candidates, edgeByChild)
	if s.questionSelector == nil {
		return fallbackQuestionSelection(candidates, &parentNode, siblings), nil
	}

	selected, err := s.questionSelector.Select(ctx, prompt, candidates)
	if err != nil {
		return fallbackQuestionSelection(candidates, &parentNode, siblings), nil
	}
	return enforceQuestionDiversity(selected, candidates, siblings), nil
}

func questionCandidates() []string {
	return []string{"なんのために？", "どうやって？", "具体的には？"}
}

func defaultQuestionCandidates() []string {
	return questionCandidates()
}

func collectAncestors(
	parentNodeID uuid.UUID,
	parentByChild map[uuid.UUID]*uuid.UUID,
	nodeByID map[uuid.UUID]model.Node,
) []model.Node {
	var ancestors []model.Node
	current := parentNodeID
	for {
		parentIDPtr, ok := parentByChild[current]
		if !ok || parentIDPtr == nil {
			break
		}
		parentID := *parentIDPtr
		if node, ok := nodeByID[parentID]; ok {
			ancestors = append(ancestors, node)
		}
		current = parentID
	}
	return ancestors
}

func collectSiblings(
	parentNodeID uuid.UUID,
	edges []model.Edge,
	nodeByID map[uuid.UUID]model.Node,
) []model.Node {
	var siblings []model.Node
	for _, edge := range edges {
		if edge.ParentNodeID == nil {
			continue
		}
		if *edge.ParentNodeID != parentNodeID {
			continue
		}
		if node, ok := nodeByID[edge.ChildNodeID]; ok {
			siblings = append(siblings, node)
		}
	}
	return siblings
}

func buildQuestionPrompt(
	parent model.Node,
	ancestors []model.Node,
	siblings []model.Node,
	candidates []string,
	edgeByChild map[uuid.UUID]model.Edge,
) string {
	var builder strings.Builder
	builder.WriteString("あなたは目標達成のための思考を促す質問を選ぶアシスタントです。\n")
	builder.WriteString("以下の候補から1つだけ選び、候補の文字列のみを出力してください。\n")
	builder.WriteString("候補が偏らないように、兄弟ノードで同じ質問が続いている場合は別の候補を優先してください。\n")
	builder.WriteString("内容が方法・手段寄りなら『具体的には？』、目的寄りなら『どうやって？』、具体例寄りなら『なんのために？』を優先してください。\n")
	builder.WriteString("候補:\n")
	for _, candidate := range candidates {
		builder.WriteString("- " + candidate + "\n")
	}
	builder.WriteString("\n")
	builder.WriteString("親ノード:\n")
	builder.WriteString(formatNodeLine(parent, edgeByChild[parent.ID]))
	builder.WriteString("\n\n")
	builder.WriteString("祖先ノード（親を除く、近い順）:\n")
	if len(ancestors) == 0 {
		builder.WriteString("- なし\n")
	} else {
		for _, node := range ancestors {
			builder.WriteString("- " + formatNodeLine(node, edgeByChild[node.ID]))
			builder.WriteString("\n")
		}
	}
	builder.WriteString("\n")
	builder.WriteString("兄弟ノード（同じ親）:\n")
	if len(siblings) == 0 {
		builder.WriteString("- なし\n")
	} else {
		for _, node := range siblings {
			builder.WriteString("- " + formatNodeLine(node, edgeByChild[node.ID]))
			builder.WriteString("\n")
		}
	}
	return builder.String()
}

func formatNodeLine(node model.Node, edge model.Edge) string {
	content := strings.TrimSpace(node.Content)
	if content == "" {
		content = "(空)"
	}
	relationInfo := ""
	if edge.Relation != "" {
		relationInfo = fmt.Sprintf(" / 関係: %s", edge.Relation)
		if edge.RelationLabel != nil && strings.TrimSpace(*edge.RelationLabel) != "" {
			relationInfo = fmt.Sprintf("%s(%s)", relationInfo, strings.TrimSpace(*edge.RelationLabel))
		}
	}
	if node.Question != nil && strings.TrimSpace(*node.Question) != "" {
		return fmt.Sprintf("内容: %s / 質問: %s%s", content, strings.TrimSpace(*node.Question), relationInfo)
	}
	return fmt.Sprintf("内容: %s%s", content, relationInfo)
}

func fallbackQuestionSelection(candidates []string, parent *model.Node, siblings []model.Node) string {
	if len(candidates) == 0 {
		return "具体的には？"
	}
	if parent != nil {
		content := strings.TrimSpace(parent.Content)
		if content != "" {
			if containsAny(content, []string{"目的", "理由", "ため", "狙い"}) {
				return pickLeastUsed(candidates, siblings, "どうやって？")
			}
			if containsAny(content, []string{"方法", "やり方", "手段", "手順", "どうやって"}) {
				return pickLeastUsed(candidates, siblings, "具体的には？")
			}
			if containsAny(content, []string{"具体", "例えば", "例", "ケース"}) {
				return pickLeastUsed(candidates, siblings, "なんのために？")
			}
		}
	}
	return pickLeastUsed(candidates, siblings, "具体的には？")
}

func enforceQuestionDiversity(selected string, candidates []string, siblings []model.Node) string {
	if selected == "" {
		return pickLeastUsed(candidates, siblings, "具体的には？")
	}
	if !isInCandidates(selected, candidates) {
		return pickLeastUsed(candidates, siblings, "具体的には？")
	}
	return pickLeastUsed(candidates, siblings, selected)
}

func pickLeastUsed(candidates []string, siblings []model.Node, preferred string) string {
	if len(candidates) == 0 {
		return "具体的には？"
	}
	counts := make(map[string]int, len(candidates))
	for _, c := range candidates {
		counts[c] = 0
	}
	for _, node := range siblings {
		if node.Question == nil {
			continue
		}
		question := strings.TrimSpace(*node.Question)
		if _, ok := counts[question]; ok {
			counts[question]++
		}
	}
	minCount := int(^uint(0) >> 1)
	for _, c := range candidates {
		if counts[c] < minCount {
			minCount = counts[c]
		}
	}
	var leastUsed []string
	for _, c := range candidates {
		if counts[c] == minCount {
			leastUsed = append(leastUsed, c)
		}
	}
	if isInCandidates(preferred, leastUsed) {
		return preferred
	}
	if len(leastUsed) > 0 {
		return leastUsed[0]
	}
	return candidates[0]
}

func isInCandidates(value string, candidates []string) bool {
	for _, c := range candidates {
		if c == value {
			return true
		}
	}
	return false
}

func containsAny(content string, keywords []string) bool {
	for _, keyword := range keywords {
		if strings.Contains(content, keyword) {
			return true
		}
	}
	return false
}
