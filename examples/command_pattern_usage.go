package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/muixstudio/flowgo"
)

func main() {
	ctx := context.Background()

	// Create process engine
	processEngine, err := flowgo.NewProcessEngineBuilder().
		WithEngineName("command-pattern-example").
		WithDatabase("postgres", "postgresql://localhost:5432/flowgo").
		WithHistory(true).
		Build()

	if err != nil {
		log.Fatalf("Failed to create process engine: %v", err)
	}

	// Start the engine
	if err := processEngine.Start(ctx); err != nil {
		log.Fatalf("Failed to start engine: %v", err)
	}
	defer processEngine.Stop(ctx)

	fmt.Println("=== FlowGo Command Pattern Example ===")
	fmt.Println()

	// Example 1: Deploy a process using DeployCommand
	fmt.Println("1. Deploying process definition using DeployCommand...")
	processJSON, err := os.ReadFile("../leave_approval.json")
	if err != nil {
		log.Fatalf("Failed to read process definition: %v", err)
	}

	deployCommand := commands.NewDeployCommand(
		"Leave Approval Process",
		"leave_approval.json",
		processJSON,
	)
	deployCommand.Category = "HR"

	// Execute command through engine
	result, err := processEngine.Execute(ctx, deployCommand)
	if err != nil {
		log.Fatalf("Deploy command failed: %v", err)
	}

	deployment := result.(*repository.Deployment)
	fmt.Printf("✓ Deployment successful: %s (ID: %s)\n", deployment.Name, deployment.ID)
	fmt.Println()

	// Example 2: Start a process instance using StartProcessInstanceCommand
	fmt.Println("2. Starting process instance using StartProcessInstanceCommand...")
	variables := map[string]interface{}{
		"applicantName": "Alice Johnson",
		"leaveType":     "Sick Leave",
		"leaveDays":     3,
		"startDate":     "2025-02-10",
		"endDate":       "2025-02-12",
		"reason":        "Medical appointment",
	}

	startCommand := commands.NewStartProcessInstanceWithBusinessKeyCommand(
		"leave-approval-process",
		"LEAVE-2025-002",
		variables,
	)

	result, err = processEngine.Execute(ctx, startCommand)
	if err != nil {
		log.Fatalf("Start process command failed: %v", err)
	}

	processInstance := result.(*runtime.ProcessInstance)
	fmt.Printf("✓ Process instance started: %s (Business Key: %s)\n",
		processInstance.ID, processInstance.BusinessKey)
	fmt.Println()

	// Example 3: Query tasks (traditional service approach)
	fmt.Println("3. Querying tasks (using service directly)...")
	taskService := processEngine.GetTaskService()
	tasks, err := taskService.CreateTaskQuery().
		ProcessInstanceID(processInstance.ID).
		Active().
		List(ctx)

	if err != nil {
		log.Fatalf("Failed to query tasks: %v", err)
	}

	if len(tasks) == 0 {
		fmt.Println("No tasks found")
	} else {
		for _, task := range tasks {
			fmt.Printf("✓ Task found: %s - %s\n", task.ID, task.Name)
		}
	}
	fmt.Println()

	// Example 4: Claim a task using ClaimTaskCommand
	if len(tasks) > 0 {
		task := tasks[0]
		fmt.Printf("4. Claiming task '%s' using ClaimTaskCommand...\n", task.Name)

		claimCommand := commands.NewClaimTaskCommand(task.ID, "alice.johnson")
		_, err = processEngine.Execute(ctx, claimCommand)
		if err != nil {
			log.Printf("Claim command failed: %v", err)
		} else {
			fmt.Printf("✓ Task claimed by alice.johnson\n")
		}
		fmt.Println()

		// Example 5: Complete the task using CompleteTaskCommand
		fmt.Printf("5. Completing task using CompleteTaskCommand...\n")
		outputVars := map[string]interface{}{
			"approved": true,
			"comment":  "Medical leave approved",
		}

		completeCommand := commands.NewCompleteTaskCommand(task.ID, outputVars)
		_, err = processEngine.Execute(ctx, completeCommand)
		if err != nil {
			log.Printf("Complete command failed: %v", err)
		} else {
			fmt.Printf("✓ Task completed successfully\n")
		}
		fmt.Println()
	}

	// Example 6: Custom command with interceptors
	fmt.Println("6. Demonstrating command interceptors...")
	fmt.Println("   (Check logs to see interceptor chain in action)")
	fmt.Println("   - LoggingInterceptor: Logs command execution")
	fmt.Println("   - TransactionInterceptor: Manages transactions")
	fmt.Println("   - ContextInterceptor: Manages CommandContext")
	fmt.Println("   - CommandInvoker: Executes the actual command")
	fmt.Println()

	fmt.Println("=== Command Pattern Benefits ===")
	fmt.Println("✓ All operations are commands - consistent execution model")
	fmt.Println("✓ Interceptor chain - logging, transactions, retry, etc.")
	fmt.Println("✓ CommandContext - shared context across command execution")
	fmt.Println("✓ Testability - easy to mock and test commands")
	fmt.Println("✓ Extensibility - easy to add new commands and interceptors")
	fmt.Println()

	fmt.Println("Example completed successfully!")
}
