package ollama

import (
	"context"
	"fmt"
	progressTool "fy-novel/internal/tools/progress"
	"fy-novel/pkg/utils"
	"io"
	"os"
	"path/filepath"
	"sync/atomic"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

const (
	ollamaImage   = "ollama/ollama"
	ollamaPort    = "11434"
	ollamaVolume  = "/ollama"
	ollamaMount   = ".ollama"
	ollamaConName = "ollama"

	// 初始化进度
	OllamaInitConTaskKey    = "ollama-init-task"
	ollamaInitConTaskAllNum = 6

	OllamaInitConSetModelTaskKey    = "ollama-set-model-task"
	ollamaInitConSetModelTaskAllNum = 6
)

var ollamaContainerID atomic.Value

var ModelTypes = map[string]bool{
	"llama3.3":            true,
	"llama3.2":            true,
	"llama3.2:1b":         true,
	"llama3.2-vision":     true,
	"llama3.2-vision:90b": true,
	"llama3.1":            true,
	"llama3.1:405b":       true,
	"phi3":                true,
	"phi3:medium":         true,
	"gemma2:2b":           true,
	"gemma2":              true,
	"gemma2:27b":          true,
	"mistral":             true,
	"moondream":           true,
	"neural-chat":         true,
	"starling-lm":         true,
	"codellama":           true,
	"llama2-uncensored":   true,
	"llava":               true,
	"solar":               true,
	"llama2":              true,
}

func GetOllamaContainerCli(ctx context.Context) (*client.Client, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("GetOllamaContainerCli failed to create Docker client: %w", err)
	}
	return cli, nil
}

func FindOllamaContainer(ctx context.Context) (bool, error) {
	cli, err := GetOllamaContainerCli(ctx)
	if err != nil {
		return false, fmt.Errorf("FindOllamaContainer failed to create Docker client: %w", err)
	}
	defer cli.Close()
	// 检查容器是否已存在
	containers, err := cli.ContainerList(ctx, container.ListOptions{
		Filters: filters.NewArgs(filters.Arg("name", ollamaConName)),
		All:     true,
	})
	if err != nil {
		return false, fmt.Errorf("FindOllamaContainer failed to list containers: %w", err)
	}

	if len(containers) > 0 {
		// 容器已存在
		con := containers[0]
		if con.State == "running" {
			ollamaContainerID.Store(con.ID)
			return true, nil
		}
		// 如果容器存在但没有运行，启动它
		if err := cli.ContainerStart(ctx, con.ID, container.StartOptions{}); err != nil {
			return false, fmt.Errorf("failed to start existing Ollama container: %w", err)
		}
		ollamaContainerID.Store(con.ID)
		return true, nil
	}
	return false, nil
}

func InitOllamaContainer(ctx context.Context) error {
	progressTool.InitTask(OllamaInitConTaskKey, ollamaInitConTaskAllNum)
	progressTool.UpdateProgress(OllamaInitConTaskKey, 3)
	defer func() {
		// 初始化完成
		progressTool.UpdateProgress(OllamaInitConTaskKey, ollamaInitConTaskAllNum)
	}()
	cli, err := GetOllamaContainerCli(ctx)
	if err != nil {
		return fmt.Errorf("InitOllamaContainer failed to create Docker client: %w", err)
	}
	defer cli.Close()

	reader, err := cli.ImagePull(ctx, ollamaImage, image.PullOptions{})
	if err != nil {
		return fmt.Errorf("InitOllamaContainer failed to pull Ollama image: %w", err)
	}
	defer reader.Close()

	// 获取用户主目录
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %w", err)
	}

	// 检查 $HOME/.ollama 文件夹是否存在
	ollamaDir := filepath.Join(homeDir, ollamaMount)
	if _, err := os.Stat(ollamaDir); os.IsNotExist(err) {
		// 如果文件夹不存在，创建它
		if err := os.MkdirAll(ollamaDir, 0755); err != nil {
			return fmt.Errorf("failed to create .ollama directory: %w", err)
		}
	}

	// 设置端口映射
	indexPort := nat.Port(fmt.Sprintf("%s/tcp", ollamaPort))
	hostConfig := &container.HostConfig{
		PortBindings: nat.PortMap{
			indexPort: []nat.PortBinding{
				{HostIP: "0.0.0.0", HostPort: ollamaPort},
			},
		},
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: ollamaDir,
				Target: ollamaVolume,
			},
		},
	}

	resp, err := cli.ContainerCreate(
		ctx,
		&container.Config{
			Image: ollamaImage,
			ExposedPorts: nat.PortSet{
				indexPort: struct{}{},
			},
		},
		hostConfig,
		nil,
		nil,
		ollamaConName,
	)
	if err != nil {
		return fmt.Errorf("failed to create Ollama container: %w", err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return fmt.Errorf("failed to start Ollama container: %w", err)
	}

	ollamaContainerID.Store(resp.ID)

	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return fmt.Errorf("failed to start Ollama container: %w", err)
		}
	case <-statusCh:
	}
	return nil
}

// GetOllamaContainerID returns the saved Ollama container ID
func GetOllamaContainerID() (string, bool) {
	id, ok := ollamaContainerID.Load().(string)
	return id, ok && id != ""
}

func OllamaContainerSetModel(ctx context.Context, modelName string) error {
	var (
		containerID string
		ok          bool
	)
	progressTool.InitTask(OllamaInitConSetModelTaskKey, ollamaInitConSetModelTaskAllNum)
	progressTool.UpdateProgress(OllamaInitConSetModelTaskKey, 2)
	defer func() {
		// 初始化完成
		progressTool.UpdateProgress(OllamaInitConTaskKey, ollamaInitConSetModelTaskAllNum)
	}()
	if hasType := ModelTypes[modelName]; !hasType {
		return fmt.Errorf("OllamaContainerSetModel unsupported model name: %s", modelName)
	}
	utils.SpinWaitMaxRetryAttempts(func() bool {
		containerID, ok = GetOllamaContainerID()
		return ok
	}, 5)
	if !ok {
		return fmt.Errorf("OllamaContainerSetModel failed to get Ollama container ID after retry")
	}

	cli, err := GetOllamaContainerCli(ctx)
	if err != nil {
		return fmt.Errorf("OllamaContainerSetModel failed to create Docker client: %w", err)
	}
	defer cli.Close()

	containerInfo, err := cli.ContainerInspect(context.Background(), containerID)
	if err != nil {
		return fmt.Errorf("OllamaContainerSetModel failed to inspect container: %w", err)
	}

	if !containerInfo.State.Running {
		return fmt.Errorf("OllamaContainerSetModel container %s is not running", containerID)
	}

	// Create exec configuration
	execConfig := container.ExecOptions{
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          []string{"ollama", "run", modelName},
	}

	// Create exec instance
	execResp, err := cli.ContainerExecCreate(context.Background(), containerID, execConfig)
	if err != nil {
		return fmt.Errorf("OllamaContainerSetModel failed to create exec: %w", err)
	}

	// Start exec instance
	execStartCheck := container.ExecAttachOptions{}
	execAttachResp, err := cli.ContainerExecAttach(
		context.Background(),
		execResp.ID,
		execStartCheck,
	)
	if err != nil {
		return fmt.Errorf("OllamaContainerSetModel failed to start exec: %w", err)
	}
	defer execAttachResp.Close()
	// Read the output from execAttachResp.Reader
	output, err := io.ReadAll(execAttachResp.Reader)
	if err != nil {
		return fmt.Errorf("OllamaContainerSetModel failed to read exec output: %w", err)
	}
	fmt.Println(output)
	return nil
}
