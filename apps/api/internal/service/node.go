package service

import (
	"context"
	"fmt"
	"hash/fnv"
	"strings"

	"github.com/google/uuid"
	"github.com/mokuhyo-driven-test/api/internal/ai"
	"github.com/mokuhyo-driven-test/api/internal/model"
	"github.com/mokuhyo-driven-test/api/internal/repository"
)

type NodeService struct {
	nodeRepo         repository.NodeRepository
	edgeRepo         repository.EdgeRepository
	questionGenerator ai.QuestionGenerator
}

func NewNodeService(nodeRepo repository.NodeRepository, edgeRepo repository.EdgeRepository, questionGenerator ai.QuestionGenerator) *NodeService {
	return &NodeService{
		nodeRepo:          nodeRepo,
		edgeRepo:          edgeRepo,
		questionGenerator: questionGenerator,
	}
}

func (s *NodeService) CreateNode(ctx context.Context, projectID uuid.UUID, req model.CreateNodeRequest) (*model.Node, *model.Edge, error) {
	var question *string
	if req.ParentNodeID != nil {
		if req.Question != nil && strings.TrimSpace(*req.Question) != "" {
			selected := strings.TrimSpace(*req.Question)
			question = &selected
		} else {
			selected, err := s.generateQuestion(ctx, projectID, *req.ParentNodeID)
			if err != nil {
				fallback := fallbackQuestionText(nil, nil, nil)
				selected = fallback
			}
			question = &selected
		}
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

func (s *NodeService) generateQuestion(ctx context.Context, projectID uuid.UUID, parentNodeID uuid.UUID) (string, error) {
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

	edgeByChild := make(map[uuid.UUID]model.Edge, len(edges))
	for _, edge := range edges {
		edgeByChild[edge.ChildNodeID] = edge
	}
	prompt := buildQuestionPrompt(parentNode, ancestors, siblings, edgeByChild)
	if s.questionGenerator == nil {
		return fallbackQuestionText(&parentNode, ancestors, siblings), nil
	}

	rawQuestion, err := s.questionGenerator.GenerateQuestion(ctx, prompt)
	if err != nil {
		return fallbackQuestionText(&parentNode, ancestors, siblings), nil
	}

	question := normalizeQuestion(rawQuestion)
	if !isQuestionUsable(question, parentNode, ancestors, siblings) {
		repairPrompt := buildQuestionRepairPrompt(rawQuestion, parentNode, ancestors, siblings, edgeByChild)
		repaired, err := s.questionGenerator.GenerateQuestion(ctx, repairPrompt)
		if err == nil {
			repairedQuestion := normalizeQuestion(repaired)
			if isQuestionUsable(repairedQuestion, parentNode, ancestors, siblings) {
				return repairedQuestion, nil
			}
		}
		return fallbackQuestionText(&parentNode, ancestors, siblings), nil
	}
	return question, nil
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
	edgeByChild map[uuid.UUID]model.Edge,
) string {
	var builder strings.Builder
	builder.WriteString("あなたは目標達成のための思考を促す質問を作るアシスタントです。\n")
	builder.WriteString("以下の条件を満たす質問を1つだけ出力してください。\n")
	builder.WriteString("- 日本語で、30字以内\n")
	builder.WriteString("- 1文で、疑問符「？」で終える\n")
	builder.WriteString("- 端的で自然な質問\n")
	builder.WriteString("- 親ノードと同じ内容を聞かない\n")
	builder.WriteString("- 兄弟ノードと同じ質問は避ける\n")
	builder.WriteString("- 親が十分具体的なら、具体化を求める質問は避ける\n")
	builder.WriteString("- 余計な説明や記号は出力しない\n\n")
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
	builder.WriteString("\n")
	builder.WriteString("兄弟ノードの質問一覧:\n")
	if len(siblings) == 0 {
		builder.WriteString("- なし\n")
	} else {
		for _, node := range siblings {
			if node.Question == nil || strings.TrimSpace(*node.Question) == "" {
				continue
			}
			builder.WriteString("- " + strings.TrimSpace(*node.Question) + "\n")
		}
	}
	return builder.String()
}

func buildQuestionRepairPrompt(
	raw string,
	parent model.Node,
	ancestors []model.Node,
	siblings []model.Node,
	edgeByChild map[uuid.UUID]model.Edge,
) string {
	var builder strings.Builder
	builder.WriteString("次の質問文を条件に合うように修正してください。\n")
	builder.WriteString("条件:\n")
	builder.WriteString("- 日本語で、30字以内\n")
	builder.WriteString("- 1文で、疑問符「？」で終える\n")
	builder.WriteString("- 親ノードと同じ内容を聞かない\n")
	builder.WriteString("- 兄弟ノードと同じ質問は避ける\n")
	builder.WriteString("- 親が十分具体的なら、具体化を求める質問は避ける\n")
	builder.WriteString("- 余計な説明や記号は出力しない\n\n")
	builder.WriteString("元の質問:\n")
	builder.WriteString(raw + "\n\n")
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
	builder.WriteString("兄弟ノードの質問一覧:\n")
	if len(siblings) == 0 {
		builder.WriteString("- なし\n")
	} else {
		for _, node := range siblings {
			if node.Question == nil || strings.TrimSpace(*node.Question) == "" {
				continue
			}
			builder.WriteString("- " + strings.TrimSpace(*node.Question) + "\n")
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

func fallbackQuestionText(parent *model.Node, ancestors []model.Node, siblings []model.Node) string {
	used := make(map[string]struct{}, len(siblings))
	for _, node := range siblings {
		if node.Question == nil {
			continue
		}
		used[strings.TrimSpace(*node.Question)] = struct{}{}
	}

	focus := "general"
	if parent != nil && parentSeemsConcrete(parent.Content) {
		focus = "purpose"
	}

	candidates := fallbackCandidatesByFocus(focus)
	filtered := filterUnused(candidates, used)
	if len(filtered) == 0 {
		filtered = candidates
	}
	if len(filtered) == 0 {
		return "その目的は何ですか？"
	}

	seed := ""
	if parent != nil {
		seed = parent.Content
	}
	index := hashIndex(seed, len(filtered))
	question := filtered[index]
	if parent != nil && parentSeemsConcrete(parent.Content) && isConcreteProbe(question) {
		for _, candidate := range filtered {
			if !isConcreteProbe(candidate) {
				return candidate
			}
		}
	}
	return question
}

func normalizeQuestion(text string) string {
	question := strings.TrimSpace(text)
	if question == "" {
		return ""
	}
	question = strings.Split(question, "\n")[0]
	question = strings.TrimSpace(question)
	question = strings.Trim(question, "「」\"'")
	question = strings.ReplaceAll(question, "?", "？")
	if !strings.HasSuffix(question, "？") {
		question += "？"
	}
	return trimToMaxRunes(question, 30)
}

func trimToMaxRunes(text string, limit int) string {
	if limit <= 0 {
		return ""
	}
	runes := []rune(text)
	if len(runes) <= limit {
		return text
	}
	trimmed := string(runes[:limit])
	if !strings.HasSuffix(trimmed, "？") {
		trimmed = strings.TrimRight(trimmed, "？")
		if len([]rune(trimmed)) >= limit {
			trimmed = string([]rune(trimmed)[:limit-1])
		}
		trimmed += "？"
	}
	return trimmed
}

func isQuestionUsable(question string, parent model.Node, ancestors []model.Node, siblings []model.Node) bool {
	trimmed := strings.TrimSpace(question)
	if trimmed == "" {
		return false
	}
	if len([]rune(trimmed)) > 30 {
		return false
	}
	if !strings.HasSuffix(trimmed, "？") {
		return false
	}
	if isDuplicateQuestion(trimmed, siblings) {
		return false
	}
	if parentSeemsConcrete(parent.Content) && isConcreteProbe(trimmed) {
		return false
	}
	if isOverlappingContent(trimmed, parent, ancestors) {
		return false
	}
	return true
}

func isDuplicateQuestion(question string, siblings []model.Node) bool {
	for _, node := range siblings {
		if node.Question == nil {
			continue
		}
		if strings.TrimSpace(*node.Question) == question {
			return true
		}
	}
	return false
}

func parentSeemsConcrete(content string) bool {
	text := strings.TrimSpace(content)
	if text == "" {
		return false
	}
	if containsAny(text, []string{"具体", "例", "例えば", "ケース", "手順", "方法", "やり方"}) {
		return true
	}
	for _, r := range text {
		if r >= '0' && r <= '9' {
			return true
		}
	}
	if containsAny(text, []string{"毎日", "毎週", "毎月", "週", "回", "時間", "分", "朝", "夜", "午前", "午後", "km", "kg"}) {
		return true
	}
	return false
}

func isConcreteProbe(question string) bool {
	return containsAny(question, []string{"具体", "詳細", "どのように", "どんな手順"})
}

func isOverlappingContent(question string, parent model.Node, ancestors []model.Node) bool {
	questionText := normalizePlainText(question)
	if questionText == "" {
		return false
	}
	content := normalizePlainText(parent.Content)
	if content == "" {
		return false
	}
	if questionText == content {
		return true
	}
	for _, node := range ancestors {
		ancestorContent := normalizePlainText(node.Content)
		if ancestorContent == "" {
			continue
		}
		if questionText == ancestorContent {
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

func normalizePlainText(text string) string {
	trimmed := strings.TrimSpace(text)
	if trimmed == "" {
		return ""
	}
	trimmed = strings.TrimSuffix(trimmed, "？")
	trimmed = strings.TrimSuffix(trimmed, "?")
	return strings.TrimSpace(trimmed)
}

func fallbackCandidatesByFocus(focus string) []string {
	switch focus {
	case "purpose":
		return []string{
			"この目標の目的は？",
			"得たい成果は何ですか？",
			"なぜそれをやりたい？",
			"成功の基準は？",
		}
	default:
		return []string{
			"最初にやる一歩は？",
			"進め方の工夫は？",
			"どこから始めますか？",
			"障害になりそうな点は？",
		}
	}
}

func filterUnused(candidates []string, used map[string]struct{}) []string {
	if len(candidates) == 0 {
		return nil
	}
	var filtered []string
	for _, candidate := range candidates {
		if _, ok := used[candidate]; ok {
			continue
		}
		filtered = append(filtered, candidate)
	}
	return filtered
}

func hashIndex(seed string, modulo int) int {
	if modulo <= 0 {
		return 0
	}
	hasher := fnv.New32a()
	_, _ = hasher.Write([]byte(seed))
	return int(hasher.Sum32()) % modulo
}
