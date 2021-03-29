package app
import (
	"bufio"
	"fmt"
	"github.com/goccy/go-yaml"
	lodestarDir "github.com/lodestar-cli/lodestar/internal/common/lodestarDir"
	"os"
)

type appEnvironment struct {
	Name string `yaml:"name"`
	SrcPath string `yaml:"srcPath"`
}

type Info struct {
	Name string   `yaml:"name"`
	Type string `yaml:"type"`
	Description string `yaml:"description"`
	RepoUrl string `yaml:"repoUrl"`
	Target string `yaml:"target"`
	StatePath string `yaml:"statePath"`
}

type LodestarAppConfig struct {
	AppInfo struct {
		Name string   `yaml:"name"`
		Type string `yaml:"type"`
		Description string `yaml:"description"`
		RepoUrl string `yaml:"repoUrl"`
		Target string `yaml:"target"`
		StatePath string `yaml:"statePath"`
	} `yaml:"appInfo"`
	EnvGraph []appEnvironment `yaml:"envGraph"`
}


//retrieves object from config tag
func GetAppConfig(path string) (*LodestarAppConfig, error){
	//read from app configuration tag
	content, err := lodestarDir.GetConfigContent(path)
	if err != nil {
		return nil, err
	}

	app, err := UnmarshalAppConfig(content)
	if err != nil {
		return nil, err
	}
	return app, nil
}

func CreateAppConfig() error{
	var info Info
	var envs []appEnvironment
	var e appEnvironment
	var correct string
	envCounter := 1

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Enter in App Info:")
	for{

		fmt.Printf("Name ['%s']: ", info.Name)
		scanner.Scan()
		h := scanner.Text()
		if h != ""{
			info.Name = h
		}
		fmt.Printf("Type ['%s']: ", info.Type)
		scanner.Scan()
		h = scanner.Text()
		if h != ""{
			info.Type = h
		}
		fmt.Printf("Description ['%s']: ", info.Description)
		scanner.Scan()
		h = scanner.Text()
		if h != ""{
			info.Description = h
		}
		fmt.Printf("App Manifest Repository ['%s']: ", info.RepoUrl)
		scanner.Scan()
		h = scanner.Text()
		if h != ""{
			info.RepoUrl = h
		}
		fmt.Printf("App Branch Target ['%s']: ", info.Target)
		scanner.Scan()
		h = scanner.Text()
		if h != ""{
			info.Target = h
		}
		fmt.Printf("Desired State Path in Manifest Repository ['%s']: ", info.StatePath)
		scanner.Scan()
		h = scanner.Text()
		if h != ""{
			info.StatePath = h
		}

		fmt.Println("-----------------")
		printAppInfo(info)
		fmt.Println("-----------------")
		fmt.Print("Does this look correct? (y/n): ")
		fmt.Scanln(&correct)

		if correct == "y"{
			break
		}
		fmt.Printf("\n[Press Enter to keep line]\n\n")
	}

	fmt.Println("Create Environment Graph:")
	for{
		fmt.Printf("Environment %d:\n", envCounter)

		fmt.Printf("Name ['%s']: ", e.Name)
		scanner.Scan()
		h := scanner.Text()
		if h != ""{
			e.Name = h
		}
		fmt.Printf("Source Path ['%s']: ", e.SrcPath)
		scanner.Scan()
		h = scanner.Text()
		if h != ""{
			e.SrcPath = h
		}
		envs = append(envs, e)

		fmt.Println("-----------------")
		printEnvironmentGraph(envs)
		fmt.Println("-----------------")
		fmt.Print("Does this look correct? (y/n): ")
		fmt.Scanln(&correct)

		if correct != "y"{
			envs = envs[:len(envs) - 1]
			fmt.Printf("\n[Press Enter to keep line]\n\n")
			continue
		}
		fmt.Print("Add another environment? (y/n): ")
		fmt.Scanln(&correct)
		if correct == "y"{
			e.SrcPath = ""
			e.Name = ""
			envCounter++
			continue
		}
		break
	}

	app := new(LodestarAppConfig)

	app.AppInfo.Name = info.Name
	app.AppInfo.Type = info.Type
	app.AppInfo.Description = info.Description
	app.AppInfo.RepoUrl = info.RepoUrl
	app.AppInfo.Target = info.StatePath
	app.EnvGraph = envs

	content, err := yaml.Marshal(app)
	if err != nil {
		return  err
	}
	fmt.Println(" Final App configuration file:")
	fmt.Println("-------------------------------")
	fmt.Println(string(content))
	lodestarDir.AddConfig(content, "app",app.AppInfo.Name)


	return nil
}

func printAppInfo(i Info) error {
	bytes, err := yaml.Marshal(i)
	if err != nil {
		return err
	}
	fmt.Print(string(bytes))

	return nil
}

func printEnvironmentGraph(e []appEnvironment) error {

	bytes, err := yaml.Marshal(e)
	if err != nil {
		return err
	}
	fmt.Println(string(bytes))

	return nil
}

func UnmarshalAppConfig(content []byte) (*LodestarAppConfig, error){
	app := LodestarAppConfig{}

	err := yaml.Unmarshal(content, &app)
	if err != nil {
		return nil, err
	}

	return &app, nil
}