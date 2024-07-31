package reposync

import (
	"bufio"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"strings"

	"github.com/go-git/go-git/v5/plumbing/object"
)

// Dependency structure for package.json
type PackageJSON struct {
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
	Keywords        []string          `json:"keywords"`
}

// Dependency structure for pom.xml
type PomXML struct {
	Dependencies []struct {
		GroupID    string `xml:"groupId"`
		ArtifactID string `xml:"artifactId"`
		Version    string `xml:"version"`
	} `xml:"dependencies>dependency"`
}

// Helper function to read lines from a file
func readLines(file *object.File) ([]string, error) {
	reader, err := file.Reader()
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	var lines []string
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

// Function to parse Pipfile
func parsePipfile(file *object.File) ([]string, error) {
	lines, err := readLines(file)
	if err != nil {
		return nil, err
	}

	dependencies := []string{}
	workingOnDependencies := false

	for _, line := range lines {
		if strings.HasPrefix(line, "[packages]") || strings.HasPrefix(line, "[dev-packages]") {
			workingOnDependencies = true
			continue
		} else if strings.Contains(line, "[") {
			workingOnDependencies = false
			continue
		}
		if workingOnDependencies {
			dep := strings.Split(line, " = ")[0]
			dependencies = append(dependencies, strings.TrimSpace(dep))
		}

	}
	return dependencies, nil
}

// Function to parse requirements.txt
func parseRequirementsTxt(file *object.File) ([]string, error) {
	lines, err := readLines(file)
	if err != nil {
		return nil, err
	}

	var dependencies []string
	for _, line := range lines {
		dep := strings.Split(line, " ")[0]
		if strings.TrimSpace(dep) != "" {
			dependencies = append(dependencies, strings.TrimSpace(dep))
		}
	}
	return dependencies, nil
}

// Function to parse package.json
func parsePackageJSON(file *object.File) ([]string, error) {
	reader, err := file.Reader()
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	// Read the file content
	content, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	var pkg PackageJSON
	if err := json.Unmarshal(content, &pkg); err != nil {
		return nil, err
	}

	var dependencies []string
	for dep := range pkg.Dependencies {
		dependencies = append(dependencies, dep)
	}
	dependencies = append(dependencies, pkg.Keywords...)

	for devDep := range pkg.DevDependencies {
		dependencies = append(dependencies, devDep)
	}
	return dependencies, nil
}

// Function to parse go.mod
func parseGoMod(file *object.File) ([]string, error) {
	lines, err := readLines(file)
	if err != nil {
		return nil, err
	}
	dependencies := []string{}
	workingOnDependencies := false
	for _, line := range lines {
		if strings.HasPrefix(line, "require") {
			workingOnDependencies = true
			continue
		} else if strings.HasPrefix(line, ")") {
			workingOnDependencies = false
			continue
		}
		if workingOnDependencies && !strings.HasSuffix(line, " // indirect") {
			dep := strings.Split(line, " ")[0]
			if strings.TrimSpace(dep) != "" {
				dependencies = append(dependencies, strings.TrimSpace(dep))
			}
		}
	}
	return dependencies, nil
}

// Function to parse pom.xml
func parsePomXML(file *object.File) ([]string, error) {

	reader, err := file.Reader()
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	// Read the file content
	content, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	var pom PomXML
	if err := xml.Unmarshal(content, &pom); err != nil {
		return nil, err
	}

	var dependencies []string
	for _, dep := range pom.Dependencies {
		dependencies = append(dependencies, fmt.Sprintf("%s:%s", dep.GroupID, dep.ArtifactID))
	}
	return dependencies, nil
}

// Function to parse build.gradle
func parseBuildGradle(file *object.File) ([]string, error) {
	lines, err := readLines(file)
	if err != nil {
		return nil, err
	}

	var dependencies []string
	inDependenciesBlock := false
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "dependencies {" {
			inDependenciesBlock = true
			continue
		}
		if inDependenciesBlock {
			if line == "}" {
				inDependenciesBlock = false
				break
			}
			componentArray := strings.Split(strings.Split(line, " '")[1], ":")
			dependencies = append(dependencies, componentArray[0]+":"+componentArray[1])
		}
	}
	return dependencies, nil
}

// Function to parse Gemfile
func parseGemfile(file *object.File) ([]string, error) {
	lines, err := readLines(file)
	if err != nil {
		return nil, err
	}

	var dependencies []string
	for _, line := range lines {
		if strings.HasPrefix(line, "gem ") {
			dep := strings.Fields(line)[1]
			dependencies = append(dependencies, strings.Trim(dep, `"',`))
		}
	}
	return dependencies, nil
}

// Function to parse Cargo.toml
func parseCargoToml(file *object.File) ([]string, error) {
	lines, err := readLines(file)
	if err != nil {
		return nil, err
	}

	var dependencies []string
	workingOnDependencies := false
	for _, line := range lines {
		if strings.HasPrefix(line, "[dependencies]") || strings.HasPrefix(line, "[dev-dependencies]") {
			workingOnDependencies = true
			continue
		} else if strings.HasPrefix(line, "[") {
			workingOnDependencies = false
			continue
		}
		if line != "" && workingOnDependencies {
			dep := strings.Split(line, "=")[0]
			dependencies = append(dependencies, strings.TrimSpace(dep))
		}
	}
	return dependencies, nil
}

// Function to parse .cabal
func parseCabal(file *object.File) ([]string, error) {
	lines, err := readLines(file)
	if err != nil {
		return nil, err
	}

	var dependencies []string
	for _, line := range lines {
		if strings.Contains(line, "build-depends:") {
			allDeps := strings.Split(line, ":")[1]
			deps := strings.Split(allDeps, ",")
			for _, dep := range deps {
				depParts := strings.Split(strings.TrimSpace(dep), " ")
				depName := strings.TrimSpace(depParts[0])
				dependencies = append(dependencies, depName)
			}

		}
	}
	return dependencies, nil
}

// Function to parse composer.json
func parseComposerJSON(file *object.File) ([]string, error) {

	reader, err := file.Reader()
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	// Read the file content
	content, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}

	var pkg map[string]interface{}
	if err := json.Unmarshal(content, &pkg); err != nil {
		return nil, err
	}

	var dependencies []string
	if require, ok := pkg["require"].(map[string]interface{}); ok {
		for dep := range require {
			dependencies = append(dependencies, dep)
		}
	}
	return dependencies, nil
}

func parseBlockChains(file *object.File) ([]string, error) {
	contents, err := file.Contents()
	if err != nil {
		return nil, err
	}
	blockChainsContained := []string{}
	lowerContents := strings.ToLower(contents)

	blockChains := []string{
		"Ethereum",
		"Gnosis",
		"Polkadot",
		"Polygon",
		"Cosmos",
		"Arbitrum",
		"BNB",
		"Avalanche",
		"Solana",
		"Optimism",
		"Bitcoin",
		"Kusama",
		"NEAR",
		"Celo",
		"Fantom",
		"Gnosis Chain",
		"Starknet",
		"Base",
		"ZKSync",
		"Moonbeam",
		"Cardano",
		"Internet Computer",
		"Aptos",
		"Moonriver",
		"IPFS",
		"Filecoin",
		"Lightning",
		"Aurora",
		"Polygon zkEVM",
		"Osmosis",
		"Tezos",
		"Mina Protocol",
		"Harmony",
		"Sui Network",
		"Oasis",
		"Hedera",
		"Stellar",
		"Status",
		"Flow",
		"Linea",
		"Acala",
		"Stacks",
		"Boba",
		"Algorand",
		"Klaytn",
		"MultiversX (Elrond)",
		"Arweave",
		"Astar Network",
		"XRP",
		"Scroll",
		"Basic Attention Token",
		"Terra",
		"EOS",
		"Fuel",
		"Chainlink",
		"Celestia",
		"TON",
		"Aztec Protocol",
		"Mantle",
		"Metis Token",
		"Litecoin",
		"Sora",
		"IOTA",
		"Cronos",
		"Monero",
		"Decentraland",
		"Injective",
		"Terra Classic",
		"Vega Protocol",
		"Sei Network",
		"The Graph",
		"EVMOS",
		"Kava.io",
		"Kujira",
		"Nostr",
		"Urbit",
		"Chia",
		"Zcash",
		"Axelar Network",
		"Anoma",
		"Holochain",
		"Huobi Token",
		"Radix DLT",
		"Kadena",
		"HECO",
		"dYdX",
		"Aragon",
		"Rootstock",
		"Audius",
		"Wormhole",
		"Dogecoin",
		"Dash",
		"Synthetix",
		"THORChain",
		"Balancer",
		"Aleo",
		"Canto",
		"Ergo",
		"Tron",
		"Golem",
		"Ocean Protocol",
	}
	for _, blockChain := range blockChains {
		lowerBlockChain := strings.ToLower(blockChain)

		if strings.Contains(lowerContents, lowerBlockChain) {
			blockChainsContained = append(blockChainsContained, lowerBlockChain)
		}
	}

	// special cases
	if strings.Contains(lowerContents, "eth-mainnet") {
		blockChainsContained = append(blockChainsContained, "ethereum")
	}

	return blockChainsContained, nil

}

func ParseFile(file *object.File, dependencyFileName string) []string {
	dependencies := []string{}
	errors := []error{}
	switch dependencyFileName {
	case "Pipfile":
		myDependencies, err := parsePipfile(file)
		dependencies = append(dependencies, myDependencies...)
		errors = append(errors, err)
	case "requirements.txt":
		myDependencies, err := parseRequirementsTxt(file)
		dependencies = append(dependencies, myDependencies...)
		errors = append(errors, err)
	case "package.json":
		myDependencies, err := parsePackageJSON(file)
		dependencies = append(dependencies, myDependencies...)
		errors = append(errors, err)
	case "go.mod":
		myDependencies, err := parseGoMod(file)
		dependencies = append(dependencies, myDependencies...)
		errors = append(errors, err)
	case "pom.xml":
		myDependencies, err := parsePomXML(file)
		dependencies = append(dependencies, myDependencies...)
		errors = append(errors, err)
	case "build.gradle":
		myDependencies, err := parseBuildGradle(file)
		dependencies = append(dependencies, myDependencies...)
		errors = append(errors, err)
	case "Gemfile":
		myDependencies, err := parseGemfile(file)
		dependencies = append(dependencies, myDependencies...)
		errors = append(errors, err)
	case "Cargo.toml":
		myDependencies, err := parseCargoToml(file)
		dependencies = append(dependencies, myDependencies...)
		errors = append(errors, err)
	case ".cabal":
		myDependencies, err := parseCabal(file)
		dependencies = append(dependencies, myDependencies...)
		errors = append(errors, err)
	case "composer.json":
		myDependencies, err := parseComposerJSON(file)
		dependencies = append(dependencies, myDependencies...)
		errors = append(errors, err)
	default:

		myDependencies, err := parseBlockChains(file)
		dependencies = append(dependencies, myDependencies...)
		errors = append(errors, err)

	}

	if errors[len(errors)-1] != nil {
		fmt.Printf("Error parsing %s: %s\n", file.Name, errors[len(errors)-1])
	}

	return getUniqueDependencies(dependencies)
}

func getUniqueDependencies(dependencies []string) []string {
	uniqueDependencies := []string{}
	dependencyMap := make(map[string]bool)
	for _, dep := range dependencies {
		if _, ok := dependencyMap[dep]; !ok {
			uniqueDependencies = append(uniqueDependencies, dep)
			dependencyMap[dep] = true
		}
	}
	return uniqueDependencies
}
