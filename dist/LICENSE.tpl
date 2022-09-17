{{.LicenseContent }}

{{- range .Groups }}
========================================================================
{{.LicenseID}} licenses
========================================================================

The following components are provided under the {{.LicenseID}} License. See project link for details.
The text of each license is also included at licenses/license-[project].txt.

{{ range .Deps }}
    {{ .Name }} {{ .Version }} {{ .LicenseID }}
{{- end }}
{{ end }}
