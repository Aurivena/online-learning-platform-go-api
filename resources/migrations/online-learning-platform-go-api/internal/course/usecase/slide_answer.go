package usecase

import (
	"context"
	"encoding/json"
	"log/slog"
	"strconv"
	"strings"

	"github.com/Aurivena/spond/v4/netsp"
	"github.com/Aurivena/spond/v4/netstatus"

	"online-learning-platform-go-api/internal/course/entity"
)

func jsonAsUint64(v interface{}) (uint64, bool) {
	if v == nil {
		return 0, false
	}
	switch x := v.(type) {
	case uint64:
		return x, true
	case uint:
		return uint64(x), true
	case int:
		if x < 0 {
			return 0, false
		}
		return uint64(x), true
	case int32:
		if x < 0 {
			return 0, false
		}
		return uint64(x), true
	case int64:
		if x < 0 {
			return 0, false
		}
		return uint64(x), true
	case float64:
		if x < 0 || x != float64(uint64(x)) {
			return 0, false
		}
		return uint64(x), true
	case json.Number:
		u, err := strconv.ParseUint(string(x), 10, 64)
		return u, err == nil
	case string:
		u, err := strconv.ParseUint(strings.TrimSpace(x), 10, 64)
		return u, err == nil
	default:
		return 0, false
	}
}

func shallowCopyPayload(p entity.PayloadJSON) entity.PayloadJSON {
	if p == nil {
		return entity.PayloadJSON{}
	}
	out := make(entity.PayloadJSON, len(p))
	for k, v := range p {
		out[k] = v
	}
	return out
}

func asPayloadMap(v interface{}) (entity.PayloadJSON, bool) {
	if v == nil {
		return nil, false
	}
	if m, ok := v.(map[string]interface{}); ok {
		return entity.PayloadJSON(m), true
	}
	if b, ok := v.([]byte); ok {
		var m map[string]interface{}
		if err := json.Unmarshal(b, &m); err == nil {
			return entity.PayloadJSON(m), true
		}
	}
	if s, ok := v.(string); ok && strings.TrimSpace(s) != "" {
		var m map[string]interface{}
		if err := json.Unmarshal([]byte(s), &m); err == nil {
			return entity.PayloadJSON(m), true
		}
	}
	return nil, false
}

func jsonTruth(v interface{}) bool {
	b, ok := v.(bool)
	return ok && b
}

// optionMarkedCorrect detects common editor flags for the right choice.
func optionMarkedCorrect(m map[string]interface{}) bool {
	for _, k := range []string{"isCorrect", "is_correct", "correct", "right", "isRight"} {
		if jsonTruth(m[k]) {
			return true
		}
	}
	switch x := m["correct"].(type) {
	case float64:
		return x != 0
	case int:
		return x != 0
	case int64:
		return x != 0
	}
	return false
}

func optionIDFromMap(m map[string]interface{}) (uint64, bool) {
	for _, key := range []string{"id", "Id", "value", "optionId", "option_id", "key"} {
		if v, ok := m[key]; ok {
			if id, ok := jsonAsUint64(v); ok {
				return id, true
			}
		}
	}
	return 0, false
}

// testOptionsBlock finds the object that holds options[] — either the root payload or a nested
// map (e.g. payload.text) matching how the LMS editor / Vue client stores TEST slides.
// Returns the block and the parent key under root ("", "text", "Text", …) for writing is_right back.
func testOptionsBlock(root entity.PayloadJSON) (block entity.PayloadJSON, parentKey string) {
	if root == nil {
		return nil, ""
	}
	if _, ok := payloadOptionsArray(root); ok {
		return root, ""
	}
	for _, key := range []string{"text", "Text", "body", "Body", "data", "Data"} {
		if m, ok := asPayloadMap(root[key]); ok {
			if _, ok := payloadOptionsArray(m); ok {
				return m, key
			}
		}
	}
	return nil, ""
}

// payloadOptionsArray decodes payload["options"] as a JSON array (handles []interface{},
// []map[string]interface{}, and other driver-specific shapes).
func payloadOptionsArray(payload entity.PayloadJSON) ([]interface{}, bool) {
	raw, ok := payload["options"]
	if !ok || raw == nil {
		return nil, false
	}
	b, err := json.Marshal(raw)
	if err != nil {
		return nil, false
	}
	var opts []interface{}
	if err := json.Unmarshal(b, &opts); err != nil {
		return nil, false
	}
	return opts, true
}

func optionExistsInBlock(block entity.PayloadJSON, optionID uint64) bool {
	opts, ok := payloadOptionsArray(block)
	if !ok {
		return false
	}
	for _, o := range opts {
		m, ok := o.(map[string]interface{})
		if !ok {
			continue
		}
		if id, ok := optionIDFromMap(m); ok && id == optionID {
			return true
		}
	}
	return false
}

// tryResolveAnswerFromLayer scans one JSON object (payload root or nested text block) for explicit
// answer ids, index fields, nested meta maps, and fuzzy key hints.
func tryResolveAnswerFromLayer(layer entity.PayloadJSON, opts []interface{}, valid map[uint64]struct{}) (uint64, bool) {
	if layer == nil {
		return 0, false
	}
	if id, ok := tryExplicitAnswerKeys(layer); ok {
		return id, true
	}
	if id, ok := tryCorrectIndex(layer, opts); ok {
		return id, true
	}
	for _, nk := range []string{"meta", "settings", "config", "quiz", "extra"} {
		if sub, ok := asPayloadMap(layer[nk]); ok {
			if id, ok := tryExplicitAnswerKeys(sub); ok {
				return id, true
			}
			if id, ok := tryCorrectIndex(sub, opts); ok {
				return id, true
			}
			if id, ok := guessCorrectByBlockKeyHints(sub, valid); ok {
				return id, true
			}
		}
	}
	if id, ok := guessCorrectByBlockKeyHints(layer, valid); ok {
		return id, true
	}
	return 0, false
}

// correctAnswerID resolves the correct option id. root is the full slide payload; block is where
// options[] was found (often payload.text). Editors often store output on the root while options
// live under text — both layers are scanned (root first, then block when nested).
func correctAnswerID(root, block entity.PayloadJSON, parentKey string, opts []interface{}) (uint64, bool) {
	valid := validOptionIDs(opts)

	layers := make([]entity.PayloadJSON, 0, 2)
	if root != nil {
		layers = append(layers, root)
	}
	if parentKey != "" && block != nil {
		layers = append(layers, block)
	}

	for _, layer := range layers {
		if id, ok := tryResolveAnswerFromLayer(layer, opts, valid); ok {
			return id, true
		}
	}

	for _, o := range opts {
		m, ok := o.(map[string]interface{})
		if !ok {
			continue
		}
		if optionMarkedCorrect(m) {
			if id, ok := optionIDFromMap(m); ok {
				return id, true
			}
		}
	}

	if len(opts) == 1 {
		if m, ok := opts[0].(map[string]interface{}); ok {
			if id, ok := optionIDFromMap(m); ok {
				return id, true
			}
		}
	}

	return 0, false
}

func validOptionIDs(opts []interface{}) map[uint64]struct{} {
	valid := make(map[uint64]struct{})
	for _, o := range opts {
		m, ok := o.(map[string]interface{})
		if !ok {
			continue
		}
		if id, ok := optionIDFromMap(m); ok {
			valid[id] = struct{}{}
		}
	}
	return valid
}

func tryExplicitAnswerKeys(b entity.PayloadJSON) (uint64, bool) {
	if b == nil {
		return 0, false
	}
	for _, key := range []string{
		"output", "correct_id", "correctId", "answer_id", "answerId",
		"correct_answer", "correctAnswer", "correct_option_id", "correctOptionId",
		"right_option_id", "rightOptionId", "solution_id", "solutionId",
		"rightAnswer", "right_answer", "true_answer", "trueAnswer",
	} {
		if v, ok := b[key]; ok {
			if id, ok := jsonAsUint64(v); ok {
				return id, true
			}
		}
	}
	return 0, false
}

func tryCorrectIndex(block entity.PayloadJSON, opts []interface{}) (uint64, bool) {
	if block == nil {
		return 0, false
	}
	for _, key := range []string{"correctIndex", "correct_index", "answerIndex", "answer_index"} {
		if v, ok := block[key]; ok {
			if idx, ok := jsonAsUint64(v); ok {
				i := int(idx)
				if i >= 0 && i < len(opts) {
					if m, ok := opts[i].(map[string]interface{}); ok {
						if id, ok := optionIDFromMap(m); ok {
							return id, true
						}
					}
				}
			}
		}
	}
	return 0, false
}

func guessCorrectByBlockKeyHints(block entity.PayloadJSON, valid map[uint64]struct{}) (uint64, bool) {
	if block == nil || len(valid) == 0 {
		return 0, false
	}
	for k, v := range block {
		lk := strings.ToLower(k)
		if lk == "options" || lk == "answers" || lk == "choices" || lk == "variants" || lk == "question" {
			continue
		}
		if !stringsContainsAny(lk, []string{"correct", "right", "solution", "answer"}) {
			continue
		}
		if id, ok := jsonAsUint64(v); ok {
			if _, ok := valid[id]; ok {
				return id, true
			}
		}
	}
	return 0, false
}

func stringsContainsAny(s string, subs []string) bool {
	for _, sub := range subs {
		if strings.Contains(s, sub) {
			return true
		}
	}
	return false
}

// CheckTestSlideOption verifies that slideID belongs to moduleID, that the slide is a TEST,
// persists payload["is_right"] and returns whether the chosen option matches payload["output"].
func (uc *SlideUseCase) CheckTestSlideOption(ctx context.Context, accountID, moduleID, slideID, optionID uint64) (bool, *netsp.Response[netsp.ErrorDetail]) {
	slog.Debug("CheckTestSlideOption", "module_id", moduleID, "slide_id", slideID, "option_id", optionID)

	mod, err := uc.moduleRepo.GetByID(ctx, moduleID)
	if err != nil {
		slog.Warn("CheckTestSlideOption: module not found", "module_id", moduleID, "err", err)
		return false, netsp.BuildError(
			netstatus.CodeNotFound,
			netsp.ErrorDetail{
				Title:    "Module Not Found",
				Message:  "The requested module does not exist",
				Solution: "Please check the module ID and try again",
			},
		)
	}
	found := false
	for i := range mod.Slides {
		if mod.Slides[i].ID == slideID {
			found = true
			break
		}
	}
	if !found {
		slog.Warn("CheckTestSlideOption: slide not in module", "module_id", moduleID, "slide_id", slideID)
		return false, netsp.BuildError(
			netstatus.CodeNotFound,
			netsp.ErrorDetail{
				Title:    "Slide Not Found",
				Message:  "This slide is not part of the given module",
				Solution: "Use a slide ID that belongs to the module",
			},
		)
	}

	slide, err := uc.slideRepo.GetByID(ctx, slideID)
	if err != nil {
		slog.Warn("CheckTestSlideOption: slide not found", "slide_id", slideID, "err", err)
		return false, netsp.BuildError(
			netstatus.CodeNotFound,
			netsp.ErrorDetail{
				Title:    "Slide Not Found",
				Message:  "The requested slide does not exist",
				Solution: "Please check the slide ID and try again",
			},
		)
	}

	if slide.SlideType != entity.SlideTypeTest {
		slog.Warn("CheckTestSlideOption: not a TEST slide", "slide_id", slideID, "slide_type", slide.SlideType)
		return false, netsp.BuildError(
			netstatus.CodeBadRequest,
			netsp.ErrorDetail{
				Title:    "Not a Test Slide",
				Message:  "Only slides of type TEST support answer checking",
				Solution: "Use a TEST slide or call the slide detail endpoint instead",
			},
		)
	}

	if slide.Payload == nil {
		slide.Payload = entity.PayloadJSON{}
	}

	block, parentKey := testOptionsBlock(slide.Payload)
	if block == nil {
		slog.Warn("CheckTestSlideOption: no options array on root or nested text/body",
			"slide_id", slideID, "payload_keys", payloadKeysForLog(slide.Payload))
		return false, netsp.BuildError(
			netstatus.CodeBadRequest,
			netsp.ErrorDetail{
				Title:    "Invalid Test Payload",
				Message:  "TEST slide has no options[] on payload root or under text/body",
				Solution: "Store options under payload.options or payload.text.options (with ids)",
			},
		)
	}

	opts, ok := payloadOptionsArray(block)
	if !ok {
		slog.Warn("CheckTestSlideOption: options key present but not decodable", "slide_id", slideID)
		return false, netsp.BuildError(
			netstatus.CodeBadRequest,
			netsp.ErrorDetail{
				Title:    "Invalid Test Payload",
				Message:  "Could not read options as a JSON array",
				Solution: "Ensure options is a JSON array of objects with id fields",
			},
		)
	}

	if !optionExistsInBlock(block, optionID) {
		slog.Warn("CheckTestSlideOption: option id not in options",
			"slide_id", slideID, "option_id", optionID, "parent_key", parentKeyOrRoot(parentKey))
		return false, netsp.BuildError(
			netstatus.CodeBadRequest,
			netsp.ErrorDetail{
				Title:    "Invalid Option",
				Message:  "The selected option is not one of this question's choices",
				Solution: "Send an option id that appears in payload.options[].id (or nested text.options)",
			},
		)
	}

	correctID, ok := correctAnswerID(slide.Payload, block, parentKey, opts)
	if !ok {
		slog.Warn("CheckTestSlideOption: cannot resolve correct answer",
			"slide_id", slideID, "parent_key", parentKeyOrRoot(parentKey),
			"root_keys", payloadKeysForLog(slide.Payload),
			"block_keys", payloadKeysForLog(block), "first_option_keys", firstOptionKeysForLog(opts))
		return false, netsp.BuildError(
			netstatus.CodeBadRequest,
			netsp.ErrorDetail{
				Title:    "Invalid Test Payload",
				Message:  "Could not find the correct option id: set payload.output (or correctIndex / isCorrect on an option), or mark one option with isCorrect",
				Solution: "Put output on the slide payload root or inside the same object as options; or use options[].isCorrect / correctIndex",
			},
		)
	}

	isRight := correctID == optionID
	if uc.resultRepo != nil {
		if err := uc.resultRepo.Upsert(ctx, accountID, moduleID, slideID, optionID, isRight); err != nil {
			slog.Error("CheckTestSlideOption: upsert result failed", "slide_id", slideID, "module_id", moduleID, "account_id", accountID, "error", err)
			return false, netsp.BuildError(
				netstatus.CodeInternalError,
				netsp.ErrorDetail{
					Title:    "Failed to Save Test Result",
					Message:  "Could not persist user test result",
					Solution: "Please retry later",
				},
			)
		}
	}

	slog.Info("CheckTestSlideOption: ok", "slide_id", slideID, "option_id", optionID, "is_right", isRight,
		"parent_key", parentKeyOrRoot(parentKey))
	return isRight, nil
}

func parentKeyOrRoot(parentKey string) string {
	if parentKey == "" {
		return "root"
	}
	return parentKey
}

// payloadKeysForLog returns top-level JSON keys for debugging (no full payload).
func payloadKeysForLog(p entity.PayloadJSON) []string {
	if p == nil {
		return nil
	}
	keys := make([]string, 0, len(p))
	for k := range p {
		keys = append(keys, k)
	}
	return keys
}

func firstOptionKeysForLog(opts []interface{}) []string {
	if len(opts) == 0 {
		return nil
	}
	m, ok := opts[0].(map[string]interface{})
	if !ok {
		return []string{"<non-object>"}
	}
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
