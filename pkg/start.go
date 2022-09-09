package pkg

// InputStruct -> Dependency Injection Data Model for StartModule
type InputStruct struct {
	Base       string
	EncodedKey string
	Owner      string
	Repo       string
}

/*
Run -> First function that needs to run to function whole pipeline
${{ github.token }} ${{ github.repository_owner }} ${{ github.event.repository.name }}

[key]

	-> {personal access token}  that needs be created on github [link](https://github.com/settings/tokens)
	-> if you are going use github action you can pass {github.token} on xxx.yaml

[owner]

	-> {owner} of the repo/project
	-> if you are going use github action you can pass {github.repository_owner} on xxx.yaml

[repo]

	-> {name} of the repository going to be used
	-> if you are going use github action you can pass {github.event.repository.name} on xxx.yaml

[base] -> It is a optional parameter. By default refers to {master} branch

[return] -> returns address of InputStruct that contains all of the function parameters
*/
func Run(key string, owner string, repo string, base ...string) (input *InputStruct) {
	if len(base) > 0 && base[0] != "" {
		return &InputStruct{Base: base[0], EncodedKey: key, Owner: owner, Repo: repo}
	}

	return &InputStruct{Base: "master", EncodedKey: key, Owner: owner, Repo: repo}
}
