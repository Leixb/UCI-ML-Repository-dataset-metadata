################################################################################
# {{ .Name }} ({{.ID}}) {{.YearCreated}}
# {{ .Abstract }}
################################################################################
@dataset {{.TaskName}} {{.Camel}} datasetdir("{{.NormSlug}}") {{if .AttrsAreAllOK}} [
    {{range .Attributes}} :{{.Name}}, {{end}}
] {{else}} false {{end}} :{{.Target}}
{{if .AttrsAreAllOK}}
preprocess(ds::{{.Name}}) = X -> coerce(X, {{range .Attributes}}
    :{{.Name}} => {{.Type}}, {{end}}
)
{{end}}

url(::{{.Name}}) = "{{.HomeURL}}"
doi(::{{.Name}}) = "{{.DOI}}"

