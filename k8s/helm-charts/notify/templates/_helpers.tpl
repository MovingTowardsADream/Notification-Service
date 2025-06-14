{{/*
Expand the name of the chart.
*/}}
{{- define "notify.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "notify.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "notify.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "notify.labels" -}}
helm.sh/chart: {{ include "notify.chart" . }}
{{ include "notify.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "notify.selectorLabels" -}}
app.kubernetes.io/name: {{ include "notify.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "notify.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "notify.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
For envirements
*/}}
{{- define "parseDotEnv" -}}
{{- $content := .Files.Get .Values.envFile }}
{{- $lines := splitList "\n" $content }}
{{- $dict := dict }}
{{- range $line := $lines }}
{{- if and $line (contains "=" $line) }}
{{- $parts := splitn "=" 2 $line }}
{{- $key := trim (index $parts._0) }}
{{- $value := trim (index $parts._1) }}
{{- if and $key (not (hasPrefix "#" $key)) }}
{{- $_ := set $dict $key $value }}
{{- end }}
{{- end }}
{{- end }}
{{- $dict | toYaml }}
{{- end }}
