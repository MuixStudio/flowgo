package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/muixstudio/flowgo"
)

func main() {
	// Create a process engine with default configuration
	engine, err := flowgo.NewProcessEngineBuilder().
		WithEngineName("flowgo-example").
		WithDatabase("postgres", "postgresql://localhost:5432/flowgo").
		WithHistory(true).
		WithAsync(true).
		Build()

	if err != nil {
		log.Fatalf("Failed to create process engine: %v", err)
	}

	// Start the engine
	ctx := context.Background()
	if err := engine.Start(ctx); err != nil {
		log.Fatalf("Failed to start process engine: %v", err)
	}
	defer engine.Stop(ctx)

	fmt.Printf("Process engine '%s' started successfully\n", engine.GetName())

	// Get services
	repositoryService := engine.GetRepositoryService()
	runtimeService := engine.GetRuntimeService()
	taskService := engine.GetTaskService()
	historyService := engine.GetHistoryService()

	// 1. Deploy a process definition
	fmt.Println("\n=== Deploying Process Definition ===")
	processDefinitionJSON, err := os.ReadFile("../leave_approval.json")
	if err != nil {
		log.Fatalf("Failed to read process definition: %v", err)
	}

	deployment, err := repositoryService.CreateDeployment().
		Name("Leave Approval Process").
		Category("HR").
		AddProcessDefinition("leave_approval.json", processDefinitionJSON).
		Deploy(ctx)

	if err != nil {
		log.Fatalf("Failed to deploy process: %v", err)
	}
	fmt.Printf("Process deployed successfully. Deployment ID: %s\n", deployment.ID)

	// 2. Query process definitions
	fmt.Println("\n=== Querying Process Definitions ===")
	processDefinitions, err := repositoryService.CreateProcessDefinitionQuery().
		ProcessDefinitionKey("leave-approval-process").
		LatestVersion().
		List(ctx)

	if err != nil {
		log.Fatalf("Failed to query process definitions: %v", err)
	}

	for _, def := range processDefinitions {
		fmt.Printf("Process Definition: %s (v%d) - %s\n", def.Key, def.Version, def.Name)
	}

	// 3. Start a process instance
	fmt.Println("\n=== Starting Process Instance ===")
	variables := map[string]interface{}{
		"applicantName": "John Doe",
		"leaveType":     "Annual Leave",
		"leaveDays":     5,
		"startDate":     "2025-02-01",
		"endDate":       "2025-02-05",
		"reason":        "Family vacation",
	}

	processInstance, err := runtimeService.StartProcessInstanceByKeyWithBusinessKey(
		ctx,
		"leave-approval-process",
		"LEAVE-2025-001",
		variables,
	)

	if err != nil {
		log.Fatalf("Failed to start process instance: %v", err)
	}
	fmt.Printf("Process instance started. ID: %s, Business Key: %s\n",
		processInstance.ID, processInstance.BusinessKey)

	// 4. Query tasks
	fmt.Println("\n=== Querying Tasks ===")
	tasks, err := taskService.CreateTaskQuery().
		ProcessInstanceID(processInstance.ID).
		Active().
		OrderByTaskCreateTime().Asc().
		List(ctx)

	if err != nil {
		log.Fatalf("Failed to query tasks: %v", err)
	}

	for _, task := range tasks {
		fmt.Printf("Task: %s - %s (Priority: %d)\n", task.ID, task.Name, task.Priority)
	}

	// 5. Claim and complete a task (if exists)
	if len(tasks) > 0 {
		task := tasks[0]
		fmt.Printf("\n=== Working on Task: %s ===", task.Name)

		// Claim the task
		if err := taskService.Claim(ctx, task.ID, "user123"); err != nil {
			log.Printf("Failed to claim task: %v", err)
		} else {
			fmt.Printf("\nTask claimed by user123")
		}

		// Add a comment
		comment, err := taskService.AddComment(ctx, task.ID, "Reviewing leave request")
		if err != nil {
			log.Printf("Failed to add comment: %v", err)
		} else {
			fmt.Printf("\nComment added: %s", comment.Message)
		}

		// Complete the task
		outputVars := map[string]interface{}{
			"approved": true,
			"comment":  "Approved - enjoy your vacation!",
		}

		if err := taskService.CompleteWithVariables(ctx, task.ID, outputVars); err != nil {
			log.Printf("Failed to complete task: %v", err)
		} else {
			fmt.Println("\nTask completed successfully")
		}
	}

	// 6. Query process variables
	fmt.Println("\n=== Process Variables ===")
	processVars, err := runtimeService.GetVariables(ctx, processInstance.ID)
	if err != nil {
		log.Fatalf("Failed to get process variables: %v", err)
	}

	for name, value := range processVars {
		fmt.Printf("%s = %v\n", name, value)
	}

	// 7. Query history
	fmt.Println("\n=== Historical Process Instances ===")
	historicProcesses, err := historyService.CreateHistoricProcessInstanceQuery().
		ProcessDefinitionKey("leave-approval-process").
		OrderByStartTime().Desc().
		List(ctx)

	if err != nil {
		log.Printf("Failed to query historic processes: %v", err)
	} else {
		for _, hist := range historicProcesses {
			status := "Running"
			if hist.EndTime != nil {
				status = "Completed"
			}
			fmt.Printf("Process: %s - Status: %s, Started: %v\n",
				hist.ID, status, hist.StartTime)
		}
	}

	fmt.Println("\n=== Example completed successfully ===")
}
