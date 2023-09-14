################################################################################
# {{ .Name }} ({{.ID}}) {{.YearCreated}}
#-------------------------------------------------------------------------------
# - {{ .NumInstances }} instances
# - {{ .NumAttributes }} attributes
#-------------------------------------------------------------------------------
# {{ .Abstract }}
################################################################################
@dataset {{.TaskName}} {{.Camel}} datasetdir("{{.NormSlug}}") {{if .AttrsAreAllOK}} [
    {{range .Attributes}} :{{.Name}}, {{end}}
] {{else}} false {{end}} :{{.Target}}
{{if and .AttrsAreAllOK .NeedsCoercion}}
preprocess(ds::{{.Name}}) = X -> coerce(X, {{range .CoerceAttrs}}
    :{{.Name}} => {{.Type}}, {{end}}
){{end}}

url(::{{.Name}}) = "{{.HomeURL}}"
doi(::{{.Name}}) = "{{.DOI}}"

