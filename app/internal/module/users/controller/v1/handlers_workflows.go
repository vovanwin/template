package usersv1

import (
	"context"
	"fmt"
	"time"

	"github.com/vovanwin/platform/pkg/logger"
	"github.com/vovanwin/template/app/internal/workflows/workflows"
	api "github.com/vovanwin/template/shared/pkg/openapi/app/v1"
	"go.temporal.io/sdk/client"
	"go.uber.org/zap"
)

// WorkflowsTestUserOnboardingPost запускает тестовый workflow пользователя
func (i Implementation) WorkflowsTestUserOnboardingPost(ctx context.Context, req *api.TestWorkflowRequest, params api.WorkflowsTestUserOnboardingPostParams) (*api.TestWorkflowResponse, error) {
	lg := logger.Named("workflows.test-user-onboarding")
	lg.Info(ctx, "Starting test user onboarding workflow",
		zap.String("user_id", req.UserID.String()),
		zap.String("email", req.Email))

	// Создаем входные данные для workflow
	workflowInput := workflows.UserOnboardingWorkflowInput{
		UserID: req.UserID.String(),
		Email:  req.Email,
		UserData: map[string]interface{}{
			"name":  req.Name,
			"email": req.Email,
		},
	}

	// Добавляем дополнительные данные если есть
	if req.AdditionalData.IsSet() {
		additionalData := req.AdditionalData.Value
		for k, v := range additionalData {
			workflowInput.UserData[k] = v
		}
	}

	// Генерируем уникальный ID для workflow
	workflowID := fmt.Sprintf("user-onboarding-%s-%d", req.UserID.String(), time.Now().Unix())

	// Запускаем workflow
	workflowRun, err := i.temporalService.GetClient().GetClient().ExecuteWorkflow(ctx, client.StartWorkflowOptions{
		ID:        workflowID,
		TaskQueue: "default-task-queue",
	}, workflows.UserOnboardingWorkflow, workflowInput)

	if err != nil {
		lg.Error(ctx, "Failed to start workflow", zap.Error(err))
		return nil, &api.ErrorStatusCode{
			StatusCode: 500,
			Response: api.Error{
				Code:    500,
				Message: "Failed to start workflow: " + err.Error(),
			},
		}
	}

	lg.Info(ctx, "Workflow started successfully",
		zap.String("workflow_id", workflowRun.GetID()),
		zap.String("run_id", workflowRun.GetRunID()))

	return &api.TestWorkflowResponse{
		Success:    true,
		WorkflowID: workflowRun.GetID(),
		RunID:      api.NewOptString(workflowRun.GetRunID()),
		Message:    fmt.Sprintf("User onboarding workflow started for user %s", req.Email),
	}, nil
}
