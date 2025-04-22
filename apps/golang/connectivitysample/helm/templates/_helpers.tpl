{{ define "tls.scheme" }}
  {{- if .Values.general.tlsEnabled -}}
https
  {{- else -}}
http
  {{- end -}}
{{ end }}

{{- define "hostAliases.external" }}
- ip: {{ .Values.general.externalIP | quote }}
  hostnames:
  - {{ .Values.general.externalHostname | quote }}
{{- end -}}

{{- define "hostAliases.vms" }}
{{- if and .Values.vms .Values.vms.ip .Values.vms.hostname }}
- ip: {{ .Values.vms.ip | quote }}
  hostnames:
  - {{ .Values.vms.hostname | quote }}
{{- end }}
{{- end -}}
