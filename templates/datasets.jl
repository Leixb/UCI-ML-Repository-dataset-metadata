################################################################################
# {{ .Name }} (ID={{.ID}}) {{.YearCreated}}
#-------------------------------------------------------------------------------
# - {{ .NumInstances }} instances
# - {{ .NumAttributes }} attributes
#-------------------------------------------------------------------------------
# {{ .Abstract }}
#-------------------------------------------------------------------------------
# # Creators
#
# {{ range .DatasetCreators }} - {{ .Creators.FirstName }} {{ .Creators.LastName }} ({{.Creators.Institution}})
# {{ end }}
#-------------------------------------------------------------------------------
# # Attribute Information
# {{range .Attributes}}{{if ne .Name " "}}
# - {{.Name}}: ({{.Type}} {{.Role}}) {{.Description}}{{end}}{{end}}
#
################################################################################
@dataset {{.TaskName}} {{.Camel}} datasetdir("{{.NormSlug}}") {{if .AttrsAreAllOK}} [
    {{range .Attributes}} :{{.Name | toCamel}}, {{end}}
] {{else}} false {{end}} :{{.Target}}
{{if and .AttrsAreAllOK .NeedsCoercion}}
preprocess(::{{.Camel}}) = X -> coerce(X, {{range .CoerceAttrs}}
    :{{.Name | toCamel}} => {{.Type | toSciType }}, {{end}}
){{end}}
{{if .HasIDAttribute}}
drop_colums(::{{.Camel}}) = [{{range .Attributes}}{{if eq .Role "ID"}}:{{.Name | toCamel}}, {{end}}{{end}}]
{{end}}
url(::{{.Camel}}) = "{{.HomeURL}}"
doi(::{{.Camel}}) = "{{.DOI}}"

