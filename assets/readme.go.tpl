# Contents

{{range .Groups -}}
{{template "RenderTOC" dict "Group" . "Lvl" 0}}
{{- end}}
{{range .Groups -}}
{{ if .IsPresentInResult }}
{{ template "RenderGroup" dict "Group" . "Lvl" 2 }}
{{ end }}
{{ end}}


{{define "RenderGroup" -}}
{{repeat "#" .Lvl}} {{.Group.SrcGroup.Title}}

{{- if .Group.SrcGroup.Description.Valid}}
_{{.Group.SrcGroup.Description.Val}}_
{{- end}}

{{range .Group.Links -}}
- [{{.Title}}]({{.SrcLink.URL}}). Stars: {{.Facts.Data.Github.StargazersCount}}
{{end}}

{{- if .Group.Groups -}}
{{$lvl := .Lvl}}
{{range .Group.Groups}}
{{template "RenderGroup" dict "Group" . "Lvl" (add $lvl 1)}}
{{end}}
{{- end }}

{{if eq .Lvl 2 -}}
[⬆ back to top](#{{anchor "Contents"}})
{{- end}}
{{- end}}


{{define "RenderTOC" -}}
{{if .Group.IsPresentInResult}}
{{- $lvl := .Lvl -}}

{{repeat "  " .Lvl}}- [{{.Group.SrcGroup.Title}}](#{{anchor .Group.SrcGroup.Title}})
{{range .Group.Groups -}}
{{template "RenderTOC" dict "Group" . "Lvl" (add $lvl 1)}}
{{- end}}
{{- end}}
{{- end}}
