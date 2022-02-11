# Commit Walker

Commit walker is a (hopefully) useful small utility that will allow you to walk all the commits in a git repository
and output a JSON file containing a structured set of data for ingestion into other tools. We use the _currently
configured branch_ in your target repository and it must be on your local disk already. At some point we could add 
the ability to clone a repository into a temporary directory and change the branch, but that's v2.

## Usage

- Clone the repository
- `go build -o walker.exe .`
- `.\walker.exe -h`
- `.\walker.exe --repository-location d:\temp\awesome-repository --output-file d:\temp\awesome.json --from-date 2022-01-01`

